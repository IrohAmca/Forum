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

func Query(id int) {
	rows, err := database.Query("SELECT * FROM people WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		var password string
		err := rows.Scan(&id,&email, &password)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Println(id,email, password)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return
	}
}


func insertData(id int, email, password string) {
	statement, err := database.Prepare("INSERT OR IGNORE INTO people (id, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, email, password)
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
