package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)
var database *sql.DB
func createDatabase() {
	database, _ = sql.Open("sqlite3", "./database.db")
	
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT, lastname TEXT)")
	statement.Exec()
}

func Query() {
	rows, err := database.Query("SELECT * FROM people")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var id int
	var name string
	var lastname string
	for rows.Next() {
		err := rows.Scan(&id, &name, &lastname)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Println("ID:", id, "Name:", name, "Lastname:", lastname)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
	}
}

func insertData() {
	statement, _ := database.Prepare("INSERT INTO people (name, lastname) VALUES (?, ?)")
	statement.Exec("John", "Doe")
}
