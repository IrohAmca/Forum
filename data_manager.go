package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func createDatabase() {
	database, _ = sql.Open("sqlite3", "./database.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, email TEXT, password TEXT)")
	statement.Exec()
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

func insertData(id int, email, password string) {
	hashedPassword := hashPassword(password)
	statement, err := database.Prepare("INSERT OR IGNORE INTO people (id, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}

func deleteData(id int) {
	statement, _ := database.Prepare("DELETE FROM people WHERE id = ?")
	statement.Exec(id)
}
