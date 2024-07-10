package db_manager

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Query_email(email string) (string, error) {
	var password string
	row := User_db.QueryRow("SELECT Email, Password FROM Users WHERE Email = ?", email)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with email %s", email)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return password, nil
}

func Query_username(email string) (string, error) {
	var username string
	row := User_db.QueryRow("SELECT Nickname FROM Users WHERE Email = ?", email)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with email %s", email)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return username, nil
}

func Query(id int) (string, string, error) {
	var email, password string
	row := User_db.QueryRow("SELECT Email, Password FROM Users WHERE UserID = ?", id)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return email, password, nil
}

func QueryToken(username string) (string, error) {
	var token string
	row := User_db.QueryRow("SELECT Token FROM Users WHERE Nickname = ?", username)
	err := row.Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with username %s", username)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return token, nil
}
func QueryTokenID(id int) (string, string, error) {
	var token, username string
	row := User_db.QueryRow("SELECT Token, Nickname FROM Users WHERE UserID = ?", id)
	err := row.Scan(&token, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return token, username, nil
}

func Query_ID(token string) (int, error) {
	var id int
	row := User_db.QueryRow("SELECT UserID FROM Users WHERE Token = ?", token)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user with token %s", token)
		}
		return 0, fmt.Errorf("error scanning row: %v", err)
	}
	return id, nil
}
func Query_ID_By_Name(username string) (int, error) {
	var id int
	row := User_db.QueryRow("SELECT UserID FROM Users WHERE Nickname = ?", username)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user with username %s", username)
		}
		return 0, fmt.Errorf("error scanning row: %v", err)
	}
	return id, nil
}
func SetMod(id int) error {
	statement, err := User_db.Prepare("UPDATE Users SET UserLevel = 1 WHERE UserID = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func SetAdmin(id int) error {
	statement, err := User_db.Prepare("UPDATE Users SET UserLevel = 2 WHERE UserID = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
func ID2Category(ids string) ([]string, error) {
	var categories []string
	dict := map[string]string{
		"1": `Gündem`,
		"2": `Ev&Yaşam`,
		"3": `Para&Ekonomi`,
		"4": `Moda&Stil`,
		"5": `İnternet&Teknoloji`,
		"6": `Eğitim&Kariyer`,
	}
	ids = strings.Replace(ids, ",", "", -1)
	ids = strings.Replace(ids, " ", "", -1)
	for _, id := range ids {
		category, ok := dict[string(id)]
		if !ok {
			return nil, fmt.Errorf("no category with id %c", id)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func Category2ID(categories []string) ([]string, error) {
	var ids []string
	dict := map[string]string{
		`Gündem`:             "1",
		`Ev&Yaşam`:           "2",
		`Para&Ekonomi`:       "3",
		`Moda&Stil`:          "4",
		`İnternet&Teknoloji`: "5",
		`Eğitim&Kariyer`:     "6",
	}
	for _, category := range categories {
		id, ok := dict[category]
		if !ok {
			return nil, fmt.Errorf("no category with name %s", category)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func GetUserName(token string) (string, error) {
	var username string
	row := User_db.QueryRow("SELECT Nickname FROM Users WHERE Token = ?", token)
	err := row.Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with token %s", token)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return username, nil
}

func CheckToken(token string) bool {
	var user string
	row := User_db.QueryRow("SELECT Token FROM Users WHERE Token = ?", token)
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	if user == "" {
		return false
	}
	return true
}
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password")
	}
	return string(hashedPassword)
}

func InsertUser(username, email, password, token string) error {
	hashedPassword := HashPassword(password)

	var exists bool
	err := User_db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE Nickname = ? OR Email = ?)", username, email).Scan(&exists)
	if err != nil {
		log.Println("Error checking for existing user:", err)
		return err
	}

	if exists {
		fmt.Println("Error: Username or email already exists.")
		return fmt.Errorf("username or email already exists")
	}

	statement, err := User_db.Prepare("INSERT INTO Users (UserLevel, Nickname, Token, Email, Password) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(0, username, token, email, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User added successfully.")
	return nil
}

func CheckTokenFromSession(token string) bool {
	var user string
	row := User_db.QueryRow("SELECT Token FROM Session WHERE Token = ?", token)
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	if user == "" {
		return false
	}
	return true
}

func InsertSession(token, cookie string) error {
	statement, err := User_db.Prepare("INSERT INTO Session (Token, Cookie) VALUES (?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(token, cookie)
	if err != nil {
		log.Println("Error executing statement", err)
		return err
	}
	return nil
}

func DeleteSession(token string) error {
	statement, err := User_db.Prepare("DELETE FROM Session WHERE Token = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(token)
	if err != nil {
		log.Println("Error executing statement:", err)
		return err
	}
	return nil
}
func GetTokenByCookie(cookie string) (string, error) {
	var token string
	row := User_db.QueryRow("SELECT Token FROM Session WHERE Cookie = ?", cookie)
	err := row.Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no session with cookie %s", cookie)
		}
		return "", fmt.Errorf("error scanning row: %v", err)
	}
	return token, nil
}

func GetTokenByName(name string) string {
	var token string
	row := User_db.QueryRow("SELECT Token FROM Users WHERE Nickname = ?", name)
	err := row.Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
	}
	return token
}

type User struct {
	ID       int
	Nickname string
	Email    string
	Password string
	Token    string
}

// GetUserByEmail function fetches the user details based on the email address.
func GetUserByEmail(email string) (*User, error) {
	var user User
	row := User_db.QueryRow("SELECT UserID, Nickname, Email, Password, Token FROM Users WHERE Email = ?", email)
	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with email %s", email)
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}
	return &user, nil
}
