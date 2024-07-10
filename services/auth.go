package services

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/db_manager"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

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
