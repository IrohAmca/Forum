package service

import (
	"errors"
	"fmt"
	"forum/models"
	"forum/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(username string) string {
	jwt_token, err := load_env("JWT_SECRET")
	if err != nil {
		fmt.Println(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwt_token))
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
	is_session := repository.CheckTokenFromSession(token)
	if is_session {
		repository.DeleteSession(token)
		c.SetCookie("cookie", "", -1, "/", "localhost", false, false)
		c.Redirect(302, "/")
	}
}
func Login(c *gin.Context) {
	var loginInput struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBind(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(loginInput.Device_Token[:len(loginInput.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	passwordReal, err := repository.Query_email(loginInput.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = Authenticate(passwordReal, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid email or password"})
		return
	}
	username, _ := repository.Query_username(loginInput.Email) // <-- This part can be optimize
	token, err := repository.QueryToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	cookie := GenerateCookie(token)

	CheckSession(token, c)
	err = repository.InsertSession(token, cookie)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome " + username, "token": token, "cookie": cookie})
}

func SignUp(c *gin.Context) {
	var user struct {
		Username     string `json:"username" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if repository.CheckDeviceToken(user.Device_Token[:len(user.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	token := GenerateToken(user.Username)
	if err := repository.InsertUser(user.Username, user.Email, user.Password, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": user.Username + " has been registered successfully"})
}

func SignOut(c *gin.Context) {
	var user struct {
		Cookie       string `json:"cookie" binding:"required"`
		Device_Token string `json:"device_token"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Reading error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(user.Device_Token[:len(user.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}
	err := repository.DeleteSession(user.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.SetCookie("cookie", "", -1, "/", "localhost", false, false)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User has been signed out successfully"})
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
	token, err := repository.GetTokenByCookie(user.Cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if repository.CheckToken(token) {
		username, err := repository.GetUserName(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid", "username": username, "token": token})
		return
	}
	if !repository.CheckToken(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token is invalid"})
		return
	}
	if repository.CheckTokenFromSession(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": true, "message": "Token is invalid"})
		return
	} // !!!
}
func ProfilePage(c *gin.Context) {
	var user_data struct {
		Username     string `json:"username" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBind(&user_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if repository.CheckDeviceToken(user_data.Device_Token[:len(user_data.Device_Token)-8]) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token"})
		return
	}

	userID, err := repository.Query_ID_By_Name(user_data.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, posts, ld_post, ld_comment, err := userInfo(userID)
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
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User data retrieved successfully", "data": data})
}

func userInfo(userID int) (*models.User, *[]models.Post, *[]models.Post, *[]models.Post, error) {
	rows, err := repository.User_db.Query("SELECT UserID, Token, UserLevel, Name, Lastname, Nickname, Email, UserBirthdate, Password FROM Users WHERE UserID = ?", userID)
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

	rows, err := repository.User_db.Query("SELECT PostID, ThreadID, Content, Likes, Dislikes, CreatedAt FROM Posts WHERE UserID = ?", userID)
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
		post.Comment, err = repository.GetCommentsByPostID(post.PostID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for i := range posts {
		rows, err = repository.User_db.Query("SELECT CategoryIDs, Title FROM Threads WHERE ThreadID = ?", posts[i].ThreadID)
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
			categories, err := repository.ID2Category(category)
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
	rows, err := repository.User_db.Query("SELECT UserID, ThreadID, Content, Likes, Dislikes, CreatedAt FROM Posts WHERE PostID = ?", PostID)
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
	rows, err = repository.User_db.Query("SELECT CategoryIDs, Title FROM Threads WHERE ThreadID = ?", post.ThreadID)
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
		categories, err := repository.ID2Category(category)
		if err != nil {
			return models.Post{}, err
		}
		post.Categories = categories
		post.Title = title
	}
	post.Comment, err = repository.GetCommentsByPostID(post.PostID)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func GetLikedDisliked(UserID int) ([]models.Post, []models.Post, error) {
	rows, err := repository.User_db.Query("SELECT PostID FROM PostLikesDislikes WHERE UserID = ?", UserID)
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
		post.Comment, err = repository.GetCommentsByPostID(postID)
		fmt.Println(post.Comment)
		if err != nil {
			return nil, nil, err
		}
		posts = append(posts, post)
	}
	rows, err = repository.User_db.Query("SELECT CommentID FROM CommentLikesDislikes WHERE UserID = ?", UserID)
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
		post.Comment, err = repository.GetCommentsByPostID(postID)
		if err != nil {
			return nil, nil, err
		}
		com_posts = append(com_posts, post)
	}

	return posts, com_posts, nil
}

func GetPostIDByCommentID(CommentID int) (int, error) {
	rows, _ := repository.User_db.Query("SELECT PostID FROM Comments WHERE CommentID = ?", CommentID)
	defer rows.Close()

	var postID int
	for rows.Next() {
		_ = rows.Scan(&postID)
	}
	fmt.Println(postID)
	return postID, nil
}
func load_env(data string) (string, error) {
	env, err := godotenv.Read("config/.env")
	if err != nil {
		return "", err
	}
	return env[data], nil
}
func CheckEmail(c *gin.Context) {
	var email struct {
		Email        string `json:"email" binding:"required"`
		Device_Token string `json:"device_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&email); err != nil {
		fmt.Println("Reading error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "user": models.User{}})
		return
	}
	if repository.CheckDeviceToken(email.Device_Token[:len(email.Device_Token)-8]) {
		fmt.Println("Invalid Device Token")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid Device Token", "user": models.User{}})
		return
	}
	user, err := repository.GetUserByEmail(email.Email)
	if err != nil {
		fmt.Println("User not found")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error(), "user": models.User{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User founded", "user": user})
}
