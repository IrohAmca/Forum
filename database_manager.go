package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var userID int
var database *sql.DB

func createDatabase() {
	database, _ = sql.Open("sqlite3", "./database.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	statement.Exec()
}

func updateUserID() {
	rows, err := database.Query("SELECT MAX(id) FROM people")
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
	row := database.QueryRow("SELECT email, password FROM people WHERE email = ?", email)
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
	row := database.QueryRow("SELECT email, password FROM people WHERE id = ?", id)
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
	statement, err := database.Prepare("INSERT OR IGNORE INTO people (id, username, email, password) VALUES (?, ?, ?, ?)")
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
	rows, err := database.Query("SELECT id, username, email, password FROM people")
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
	statement, _ := database.Prepare("DELETE FROM people WHERE id = ?")
	statement.Exec(id)
}*/
