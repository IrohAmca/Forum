package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"frontend/manager"
	"frontend/models"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func checkEmail(email string) (models.User, error) {
	var loginInput struct {
		Email        string `json:"email"`
		Device_Token string `json:"device_token"`
	}
	env, err := godotenv.Read("config/.env")
	if err != nil {
		return models.User{}, err
	}
	loginInput.Device_Token = env["DEVICE_TOKEN"]
	loginInput.Email = email
	loginData, err := json.Marshal(loginInput)
	if err != nil {
		return models.User{}, err
	}
	api := manager.API{}
	url := api.GetURL("CheckEmail")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		fmt.Println("Error: ", err)
		return models.User{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return models.User{}, err
	}
	var response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		User    models.User `json:"user"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Marshal error: ", err)
		return models.User{}, err
	}
	success := response.Success
	if !success {
		fmt.Println("Response message: ", response.Message)
		return models.User{}, fmt.Errorf(response.Message)
	}
	return response.User, nil
}
func authSignUp(email, name string) error {
	var loginInput struct {
		Email        string `json:"email"`
		Username     string `json:"username"`
		Device_Token string `json:"device_token"`
	}
	loginInput.Email = email
	loginInput.Username = name
	env, err := godotenv.Read("config/.env")
	if err != nil {
		return err
	}
	loginInput.Device_Token = env["DEVICE_TOKEN"]
	loginData, err := json.Marshal(loginInput)
	if err != nil {
		return err
	}
	api := manager.API{}
	url := api.GetURL("AuthSignup")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf(response.Message)
	}
	return nil
}
func authLogin(email string) (string, error) {
	var loginInput struct {
		Email        string `json:"email"`
		Device_Token string `json:"device_token"`
	}
	loginInput.Email = email
	env, err := godotenv.Read("config/.env")
	if err != nil {
		return "", err
	}
	loginInput.Device_Token = env["DEVICE_TOKEN"]
	loginData, err := json.Marshal(loginInput)
	if err != nil {
		return "", err
	}
	api := manager.API{}
	url := api.GetURL("AuthLogin")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(loginData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Cookie  string `json:"cookie"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if !response.Success {
		return "", fmt.Errorf(response.Message)
	}
	return response.Cookie, nil
}
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
		fmt.Printf("\n oauthConf.Exchange() failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		fmt.Println("Failed to get user info")
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
		fmt.Println("Failed to decode JSON response")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	existingUser, err := checkEmail(googleUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", googleUser.Email) {
			err := authSignUp(googleUser.Email, googleUser.Name)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			cookie, err := authLogin(googleUser.Email)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", existingUser.Token.String)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			fmt.Println("Error: ", err)
			return
		}
	}

	new_cookie, err := authLogin(googleUser.Email)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	c.SetCookie("cookie", new_cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token.String)
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

	existingUser, err := checkEmail(githubUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", githubUser.Email) {
			err := authSignUp(githubUser.Email, githubUser.Login)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			cookie, err := authLogin(githubUser.Email)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", existingUser.Token.String)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			fmt.Println("Error: ", err)
			return
		}
	}

	new_cookie, err := authLogin(githubUser.Email)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	c.SetCookie("cookie", new_cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token.String)
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

	existingUser, err := checkEmail(fbUser.Email)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user with email %s", fbUser.Email) {
			err := authSignUp(fbUser.Email, fbUser.Name)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			cookie, err := authLogin(fbUser.Email)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			c.SetCookie("cookie", cookie, 3600, "/", "localhost", false, false)
			c.Header("Authorization", existingUser.Token.String)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			fmt.Println("Error: ", err)
			return
		}
	}

	new_cookie, err := authLogin(fbUser.Email)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	c.SetCookie("cookie", new_cookie, 3600, "/", "localhost", false, false)
	c.Header("Authorization", existingUser.Token.String)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
