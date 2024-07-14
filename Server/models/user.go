package models

import "database/sql"

type User struct {
	UserID        int
	Token         sql.NullString
	UserLevel     int
	Name          sql.NullString
	Lastname      sql.NullString
	Nickname      string
	Email         string
	UserBirthdate sql.NullString
	Password      sql.NullString
}
