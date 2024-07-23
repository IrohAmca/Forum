package models

type Comment struct {
	CommentID      int
	UserID         int
	PostID         int
	Content        string
	Username       string
	CreatedAt      string
	LikeCounter    int
	DislikeCounter int
	UpdatedAt      string
}
