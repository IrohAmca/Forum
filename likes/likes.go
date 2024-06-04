package likes

import (
	"database/sql"
	"fmt"
)

type LikeDislikeActions struct {
	UserID int  // user that performing the action.
	PostID int  // id of post that being like or dislike.
	IsLike bool // true for like, false for dislike.
}

func HandleLikeDislike(user_db *sql.DB, action LikeDislikeActions) error {

	var currentID int
	var currentIsLike bool

	//check if the user has already liked or disliked this post.
	err := user_db.QueryRow(`SELECT ID, IsLike FROM LikesDislikes WHERE UserID = ? AND PostID = ?`, action.UserID, action.PostID).Scan(&currentID, &currentIsLike)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("query error: %v", err)
	}

	if err == sql.ErrNoRows {
		// if there is no record, insert a new record.
		_, err = user_db.Exec(`INSERT INTO LikesDislikes (UserID, PostID, IsLike) VALUES (?,?,?) `, action.UserID, action.PostID, action.IsLike)
		if err != nil {
			return fmt.Errorf("insert error: %v", err)
		}

		if action.IsLike {
			_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes +1 WHERE PostID = ?`, action.PostID)

		} else {
			_, err = user_db.Exec(`UPDATE Posts SET Dislike = Dislike +1 WHERE PostID = ?`, action.PostID)
		}
		if err != nil {
			return fmt.Errorf("update count error: %v", err)
		}

	} else {

		if currentIsLike == action.IsLike {
			_, err = user_db.Exec(`DELETE FROM  LikesDislikes WHERE ID`, currentID)
			if err != nil {
				return fmt.Errorf("delete error: %v", err)
			}
			if action.IsLike {
				_, err = user_db.Exec(`UPDATE Post SET Likes = Likes -1 WHERE PostID =?`, action.PostID)
			} else {
				_, err = user_db.Exec(`UPDATE Post SET Dislike = Dislike -1 WHERE PostID =?`, action.PostID)
			}

			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		} else {
			_, err = user_db.Exec(`UPDATE LikesDislikes SET IsLike = ? WHERE ID =?`, action.IsLike, currentID)
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}

			if action.IsLike {
				_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes +1, Dislike = Dislike-1 WHERE ID =?`, action.PostID)
			} else {
				_, err = user_db.Exec(`UPDATE Posts SET Likes = Likes -1, Dislike = Dislike +1 WHERE ID =?`, action.PostID)
			}
			if err != nil {
				return fmt.Errorf("update count error: %v", err)
			}
		}

	}
	return nil
}
