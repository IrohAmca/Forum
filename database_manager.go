package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var userID int
var user_db *sql.DB

func createDatabase() {
	user_db, _ = sql.Open("sqlite3", "./user.db")
	usr_statement, _ := user_db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	usr_statement.Exec()

	post_db, _ := sql.Open("sqlite3", "./post.db")
	post_statement, _ := post_db.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT, category TEXT, date TEXT)")
	post_statement.Exec()

	post_db.Exec("ALTER TABLE posts ADD COLUMN user_id INTEGER")
	post_db.Exec("UPDATE posts SET user_id = ?", userID)

}

func updateUserID() {
	rows, err := user_db.Query("SELECT MAX(id) FROM people")
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	var maxID int
	if rows.Next() {
		err := rows.Scan(&maxID)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
	}

	userID = maxID + 1
}

func Query_email(email string) (string, string, error) {
	var email_, password string
	row := user_db.QueryRow("SELECT email, password FROM people WHERE email = ?", email)
	err := row.Scan(&email_, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with email %s", email)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return email_, password, nil
}

func Query(id int) (string, string, error) {
	var email, password string
	row := user_db.QueryRow("SELECT email, password FROM people WHERE id = ?", id)
	err := row.Scan(&email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no user with id %d", id)
		}
		return "", "", fmt.Errorf("error scanning row: %v", err)
	}
	return email, password, nil
}

func insertData(id int, username, email, password string) {
	hashedPassword := hashPassword(password)
	statement, err := user_db.Prepare("INSERT OR IGNORE INTO people (id, username, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, username, email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}

func WriteAllData() {
	rows, err := user_db.Query("SELECT id, username, email, password FROM people")
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email, password, username string
		err := rows.Scan(&id, &username, &email, &password)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("ID: %d, username: %s, Email: %s, Password: %s:\n", id, username, email, password)
	}
}

/*
func deleteData(id int) {
	statement, _ := user_db.Prepare("DELETE FROM people WHERE id = ?")
	statement.Exec(id)
}*/
