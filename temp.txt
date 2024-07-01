package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/gologin/google"
	"github.com/dghubble/sessions"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

var (
	db           *sql.DB
	sessionStore = sessions.NewCookieStore([]byte("secret-key"), nil)
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./user.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create user table if it doesn't exist
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT UNIQUE,
        hashed_password TEXT,
        provider TEXT,
        provider_id TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func main() {
	// Google OAuth2 config
	googleOAuth2Config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	}

	// GitHub OAuth2 config
	githubOAuth2Config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}

	http.Handle("/auth/google", google.StateHandler(gologin.DebugOnlyCookieConfig, google.LoginHandler(googleOAuth2Config, nil)))
	http.Handle("/auth/google/callback", google.StateHandler(gologin.DebugOnlyCookieConfig, google.CallbackHandler(googleOAuth2Config, issueSession(), nil)))

	http.Handle("/auth/github", github.StateHandler(gologin.DebugOnlyCookieConfig, github.LoginHandler(githubOAuth2Config, nil)))
	http.Handle("/auth/github/callback", github.StateHandler(gologin.DebugOnlyCookieConfig, github.CallbackHandler(githubOAuth2Config, issueSession(), nil)))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		token, err := oauth2.TokenFromContext(ctx)
		if err != nil {
			http.Error(w, "Token Error", http.StatusInternalServerError)
			return
		}
		userInfo, err := fetchUserInfo(token)
		if err != nil {
			http.Error(w, "Error fetching user info", http.StatusInternalServerError)
			return
		}

		// Generate UUID for the user
		userID := uuid.New().String()

		// Save user to database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("defaultpassword"), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`INSERT OR IGNORE INTO users (id, email, hashed_password, provider, provider_id) VALUES (?, ?, ?, ?, ?)`,
			userID, userInfo.Email, string(hashedPassword), userInfo.Provider, userInfo.ProviderID)
		if err != nil {
			http.Error(w, "Error saving user to database", http.StatusInternalServerError)
			return
		}

		session := sessionStore.New("example-session")
		session.Values["access_token"] = token.AccessToken
		session.Save(w)
		fmt.Fprintf(w, "Login successful, token: %s", token.AccessToken)
	}
	return http.HandlerFunc(fn)
}

func fetchUserInfo(token *oauth2.Token) (*UserInfo, error) {
	// Fetch user info from Google or GitHub using the access token
	// Implement this function based on the OAuth provider you are using
	return &UserInfo{}, nil
}

type UserInfo struct {
	Email      string
	Provider   string
	ProviderID string
}
