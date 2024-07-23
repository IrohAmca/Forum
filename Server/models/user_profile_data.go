package models

type UserData struct {
	User     *User
	Posts    *[]Post
	LD_Posts *[]Post
	LD_Comment *[]Post
}
