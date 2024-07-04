package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum/db_manager"
	"forum/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func GenerateToken(username string) string {
	load_env()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic("Failed to generate token")
	}
	return tokenString
}

func GenerateCookie(token string) string {
	cookie := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": token,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	cookieString, err := cookie.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic("Failed to generate token")
	}
	return cookieString
}

func CheckSession(token string, c *gin.Context) {
	is_session := db_manager.CheckTokenFromSession(token)
	if is_session {
		db_manager.DeleteSession(token)
		c.SetCookie("cookie", "", -1, "/", "localhost", false, false)
		c.Redirect(302, "/")
	}
}
func Login(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	passwordReal, err := db_manager.Query_email(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = Authenticate(passwordReal, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}
	username, _ := db_manager.Query_username(loginInput.Email) // <-- This part can be optimize
	token, err := db_manager.QueryToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	cookie := GenerateCookie(token)

	CheckSession(token, c)
	err = db_manager.InsertSession(token, cookie)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", token)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + username, "token": token})
}

func SignUp(c *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	token := GenerateToken(user.Username)
	if err := db_manager.InsertUser(user.Username, user.Email, user.Password, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": user.Username + " has been registered successfully"})
}

func SignOut(c *gin.Context) {
	c.SetCookie("cookie", "", -1, "/", "localhost", false, false)
	c.JSON(200, gin.H{"success": true, "message": "You have been signed out"})
	c.Redirect(302, "/")
	db_manager.DeleteSession(c.GetHeader("Authorization"))
}

func Authenticate(storedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		fmt.Println("Password does not match") // add new error message
		return errors.New("authentication failed")
	}
	return nil
}

func UserChecker(c *gin.Context) {
	var user struct {
		Cookie string `json:"cookie" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	token, err := db_manager.GetTokenByCookie(user.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if db_manager.CheckToken(token) {
		username, err := db_manager.GetUserName(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid", "username": username, "token": token})
		return
	}
	if !db_manager.CheckToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token is invalid"})
		return
	}
	if db_manager.CheckTokenFromSession(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": true, "message": "Token is invalid"})
		return
	}
}
func ProfilePage(c *gin.Context) {
	username := c.Param("username")
	w := c.Writer
	userID, err := db_manager.Query_ID_By_Name(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, posts, ld_post, ld_comment, err := userInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	tpl, err := template.ParseFiles("templates/userprofile.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	data := models.UserData{
		User:       user,
		Posts:      posts,
		LD_Posts:   ld_post,
		LD_Comment: ld_comment,
	}
	fmt.Println(data.Posts)
	tpl.Execute(w, &data)
}

func userInfo(userID int) (*models.User, *[]models.Post, *[]models.Post, *[]models.Post, error) {
	rows, err := db_manager.User_db.Query("SELECT UserID, Token, UserLevel, Name, Lastname, Nickname, Email, UserBirthdate, Password FROM Users WHERE UserID = ?", userID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var user models.User
	if rows.Next() {
		err := rows.Scan(
			&user.UserID,
			&user.Token,
			&user.UserLevel,
			&user.Name,
			&user.Lastname,
			&user.Nickname,
			&user.Email,
			&user.UserBirthdate,
			&user.Password,
		)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		return nil, nil, nil, nil, errors.New("user not found")
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}
	Posts := []models.Post{}
	LD_Posts := []models.Post{}
	LD_Comment := []models.Post{}
	LD_Posts, LD_Comment, err = GetLikedDisliked(user.UserID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	Posts, err = fetchPostsByUserID(user.UserID)

	if err != nil {
		return nil, nil, nil, nil, err
	}
	return &user, &Posts, &LD_Posts, &LD_Comment, err
}

func fetchPostsByUserID(userID int) ([]models.Post, error) {

	rows, err := db_manager.User_db.Query("SELECT PostID, ThreadID, Content, Likes, Dislikes, CreatedAt FROM Posts WHERE UserID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() { // <--- This part can be optimize
		var post models.Post
		err := rows.Scan(
			&post.PostID,
			&post.ThreadID,
			&post.Content,
			&post.LikeCounter,
			&post.DislikeCounter,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		post.Comment, err = db_manager.GetCommentsByPostID(post.PostID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for i := range posts {
		rows, err = db_manager.User_db.Query("SELECT CategoryIDs, Title FROM Threads WHERE ThreadID = ?", posts[i].ThreadID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var category string
			var title string
			err := rows.Scan(
				&category,
				&title,
			)
			if err != nil {
				return nil, err
			}
			categories, err := db_manager.ID2Category(category)
			if err != nil {
				return nil, err
			}
			posts[i].Categories = categories
			posts[i].Title = title
		}
	}
	return posts, nil
}

func fetchPostsByPostID(PostID int) (models.Post, error) {
	rows, err := db_manager.User_db.Query("SELECT UserID, ThreadID, Content, Likes, Dislikes, CreatedAt FROM Posts WHERE PostID = ?", PostID)
	if err != nil {
		return models.Post{}, err
	}
	defer rows.Close()
	var post models.Post
	for rows.Next() {
		err := rows.Scan(
			&post.UserID,
			&post.ThreadID,
			&post.Content,
			&post.LikeCounter,
			&post.DislikeCounter,
			&post.CreatedAt,
		)
		if err != nil {
			return models.Post{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return models.Post{}, err
	}
	rows, err = db_manager.User_db.Query("SELECT CategoryIDs, Title FROM Threads WHERE ThreadID = ?", post.ThreadID)
	if err != nil {
		return models.Post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var category string
		var title string
		err := rows.Scan(
			&category,
			&title,
		)
		if err != nil {
			return models.Post{}, err
		}
		categories, err := db_manager.ID2Category(category)
		if err != nil {
			return models.Post{}, err
		}
		post.Categories = categories
		post.Title = title
	}
	post.Comment, err = db_manager.GetCommentsByPostID(post.PostID)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func GetLikedDisliked(UserID int) ([]models.Post, []models.Post, error) {
	rows, err := db_manager.User_db.Query("SELECT PostID FROM PostLikesDislikes WHERE UserID = ?", UserID)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	var posts []models.Post
	var com_posts []models.Post
	for rows.Next() {
		var post models.Post
		var postID int
		_ = rows.Scan(&postID)
		post, _ = fetchPostsByPostID(postID)
		post.Comment, err = db_manager.GetCommentsByPostID(postID)
		fmt.Println(post.Comment)
		if err != nil {
			return nil, nil, err
		}
		posts = append(posts, post)
	}
	rows, err = db_manager.User_db.Query("SELECT CommentID FROM CommentLikesDislikes WHERE UserID = ?", UserID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		var commentID int
		_ = rows.Scan(&commentID)
		postID, err := GetPostIDByCommentID(commentID)
		if err != nil {
			return nil, nil, err
		}
		post, _ = fetchPostsByPostID(postID)
		post.Comment, err = db_manager.GetCommentsByPostID(postID)
		if err != nil {
			return nil, nil, err
		}
		com_posts = append(com_posts, post)
	}

	return posts, com_posts, nil
}

func GetPostIDByCommentID(CommentID int) (int, error) {
	rows, _ := db_manager.User_db.Query("SELECT PostID FROM Comments WHERE CommentID = ?", CommentID)
	defer rows.Close()

	var postID int
	for rows.Next() {
		_ = rows.Scan(&postID)
	}
	fmt.Println(postID)
	return postID, nil
}
func load_env() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     "GOOGLE_CLIENT_ID",
		ClientSecret: "GOOGLE_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "random"
)

func GoogleLogin(c *gin.Context) {
	load_env()
	var ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	var ClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	oauthConf.ClientID = ClientID
	oauthConf.ClientSecret = ClientSecret
	oauthStateString = "hello world"
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	if c.Request.FormValue("state") != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Request.FormValue("code")
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Mevcut kullanıcıyı email ile kontrol et
	existingUser, err := db_manager.GetUserByEmail(googleUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", googleUser.Email) {
			// Yeni kullanıcı ekle
			user_token := GenerateToken(googleUser.Name)
			err := db_manager.InsertUser(googleUser.Name, googleUser.Email, "", user_token)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}

			cookie := GenerateCookie(user_token)
			err = db_manager.InsertSession(user_token, cookie)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", user_token)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	}

	// Kullanıcı zaten mevcutsa, kullanıcı bilgileriyle giriş yap
	cookie := GenerateCookie(existingUser.Token)
	err = db_manager.InsertSession(existingUser.Token, cookie)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

var (
	githuboauthConf = &oauth2.Config{
		ClientID:     "GITHUB_CLIENT_ID",
		ClientSecret: "GITHUB_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: github.Endpoint,
	}
	githuboauthStateString = "random"
)

func GithubLogin(c *gin.Context) {
	load_env()
	var ClientID = os.Getenv("GITHUB_CLIENT_ID")
	var ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	githuboauthConf.ClientID = ClientID
	githuboauthConf.ClientSecret = ClientSecret
	githuboauthStateString = "hello world"
	url := githuboauthConf.AuthCodeURL(githuboauthStateString, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)

}

func GithubCallback(c *gin.Context) {
	if c.Request.FormValue("state") != githuboauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", githuboauthStateString, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Request.FormValue("code")
	token, err := githuboauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	client := githuboauthConf.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Printf("client.Get() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		fmt.Printf("json.NewDecoder() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Eğer email boşsa email endpoint'ini sorgula
	if githubUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			fmt.Printf("client.Get() for emails failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		defer emailResp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			fmt.Printf("json.NewDecoder() for emails failed with '%s'\n", err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Primary ve verified email adresini bul
		for _, email := range emails {
			if email.Primary && email.Verified {
				githubUser.Email = email.Email
				break
			}
		}
	}

	if githubUser.Email == "" {
		fmt.Println("No verified primary email found for the GitHub user.")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Mevcut kullanıcıyı email ile kontrol et
	existingUser, err := db_manager.GetUserByEmail(githubUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", githubUser.Email) {
			// Yeni kullanıcı ekle
			user_token := GenerateToken(githubUser.Login)
			err := db_manager.InsertUser(githubUser.Login, githubUser.Email, "", user_token)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}

			cookie := GenerateCookie(user_token)
			err = db_manager.InsertSession(user_token, cookie)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", user_token)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	}

	// Kullanıcı zaten mevcutsa, kullanıcı bilgileriyle giriş yap
	cookie := GenerateCookie(existingUser.Token)
	err = db_manager.InsertSession(existingUser.Token, cookie)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

var (
	facebookoauthConf = &oauth2.Config{
		ClientID:     "FACEBOOK_CLIENT_ID",
		ClientSecret: "FACEBOOK_CLIENT_SECRET",
		RedirectURL:  "http://localhost:8080/auth/facebook/callback",
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	facebookoauthStateString = "random"
)

func FacebookLogin(c *gin.Context) {
	load_env()
	var ClientID = os.Getenv("FACEBOOK_CLIENT_ID")
	var ClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	facebookoauthConf.ClientID = ClientID
	facebookoauthConf.ClientSecret = ClientSecret
	facebookoauthStateString = "hello world"
	url := facebookoauthConf.AuthCodeURL(facebookoauthStateString, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func FacebookCallback(c *gin.Context) {
	if c.Request.FormValue("state") != facebookoauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", facebookoauthStateString, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Request.FormValue("code")
	token, err := facebookoauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=id,name,email", token.AccessToken))
	if err != nil {
		fmt.Printf("http.Get() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer resp.Body.Close()

	var fbUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fbUser); err != nil {
		fmt.Printf("json.NewDecoder() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// Mevcut kullanıcıyı email ile kontrol et
	existingUser, err := db_manager.GetUserByEmail(fbUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", fbUser.Email) {
			// Yeni kullanıcı ekle
			user_token := GenerateToken(fbUser.Name)
			err := db_manager.InsertUser(fbUser.Name, fbUser.Email, "", user_token)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}

			cookie := GenerateCookie(user_token)
			err = db_manager.InsertSession(user_token, cookie)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", user_token)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
	}

	// Kullanıcı zaten mevcutsa, kullanıcı bilgileriyle giriş yap
	cookie := GenerateCookie(existingUser.Token)
	err = db_manager.InsertSession(existingUser.Token, cookie)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
