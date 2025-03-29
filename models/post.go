package models

import (
	"database/sql"
	"time"
)

// Post represents a post in the system
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatePost creates a new post in the database
func CreatePost(userID int, title, content string) (int, error) {
	query := "INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)"
	result, err := db.Exec(query, userID, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetPostByID retrieves a post by its ID
func GetPostByID(id int) (*Post, error) {
	query := "SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?"
	row := db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}

// GetPostsByUserID retrieves all posts for a user
func GetPostsByUserID(userID int) ([]*Post, error) {
	query := "SELECT id, user_id, title, content, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetAllPosts retrieves all posts
func GetAllPosts() ([]*Post, error) {
	query := "SELECT id, user_id, title, content, created_at FROM posts ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// UpdatePost updates a post
func UpdatePost(id int, title, content string) error {
	query := "UPDATE posts SET title = ?, content = ? WHERE id = ?"
	_, err := db.Exec(query, title, content, id)
	return err
}

// DeletePost deletes a post
func DeletePost(id int) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}
