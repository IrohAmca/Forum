package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var userID = 1
var database *sql.DB

func createDatabase() {
	database, _ = sql.Open("sqlite3", "./database.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, email TEXT, password TEXT, age INTEGER)")
	statement.Exec()
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

func insertData(id int, email, password string, age int) {
	hashedPassword := hashPassword(password)
	statement, err := database.Prepare("INSERT OR IGNORE INTO people (id, email, password, age) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, email, hashedPassword, age)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}
func WriteAllData() {
	rows, err := database.Query("SELECT id, email, password, age FROM people")
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, age int
		var email, password string
		err := rows.Scan(&id, &email, &password, &age)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("ID: %d, Email: %s, Password: %s, Age: %d\n", id, email, password, age)
	}
}

/*
func deleteData(id int) {
	statement, _ := database.Prepare("DELETE FROM people WHERE id = ?")
	statement.Exec(id)
}*/
