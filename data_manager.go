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

func Query(id int) {
	rows, err := database.Query("SELECT * FROM people WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var lastname string
		err := rows.Scan(&id,&name, &lastname)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Println(id,name, lastname)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return
	}
}


func insertData(id int, name, lastname string) {
	statement, err := database.Prepare("INSERT OR IGNORE INTO people (id, name, lastname) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, name, lastname)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	fmt.Println("Data inserted successfully.")
}

func deleteData(data string) {
	statement, _ := database.Prepare("DELETE FROM people WHERE name = ?")
	statement.Exec(data)
}
