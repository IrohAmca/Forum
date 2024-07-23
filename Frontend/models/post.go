package models

type Post struct {
	PostID         int
	ThreadID       int
	Title          string
	UserToken      string
	UserID         int
	Username       string
	Content        string
	Image          string
	Categories     []string
	CreatedAt      string
	LikeCounter    int
	DislikeCounter int
	Comment        []Comment
}
