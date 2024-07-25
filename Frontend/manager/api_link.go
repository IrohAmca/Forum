package manager

import (
	"os"

	"github.com/joho/godotenv"
)

type API struct {
	ProfilePage        string
	GoogleLogin        string
	GoogleCallback     string
	GithubLogin        string
	GithubCallback     string
	FacebookLogin      string
	FacebookCallback   string
	SignUp             string
	Login              string
	UserChecker        string
	SignOut            string
	GetPosts           string
	LikeDislikeComment string
	DeleteComment      string
	LikeDislikePost    string
	CreatePost         string
	DeletePost         string
	CreateComment      string
}

var dictionary = map[string]string{
	"ProfilePage":        "/profile",             // OK!!!
	"AuthLogin":          "/auth/login",          // OK!!!
	"AuthSignup":         "/auth/signup",		  // OK!!!
	"SignUp":             "/sign-up",			  // OK!!!
	"Login":              "/login",               // OK!!!
	"UserChecker":        "/check-token",         // OK!!!
	"SignOut":            "/sign-out",            // OK!!!
	"GetPosts":           "/get-posts",           // OK!!!
	"LikeDislikeComment": "/ld_comment",          // OK!!!
	"DeleteComment":      "/delete-comment",      // OK!!!
	"LikeDislikePost":    "/ld_post",			  // OK!!!
	"CreatePost":         "/create-post",         // OK!!!
	"DeletePost":         "/delete-post",         // OK!!!
	"CreateComment":      "/create-comment",	  // OK!!!
	"DeviceRegister":     "/device-register",     // OK!!!
	"CheckEmail":         "/check-email",         // OK!!!
	"SetModarator":       "/setModarator",        
}

func init() {
	err := godotenv.Load("config/.env")
	if err != nil {
		panic(err)
	}
}

func (a *API) GetURL(endpointName string) string {
	url := os.Getenv("API_URL")
	return "http://" + url + dictionary[endpointName]
}
