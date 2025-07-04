// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package repository

import (
	"database/sql"
	"time"
)

type Link struct {
	ID           int64          `json:"id"`
	Url          string         `json:"url"`
	Title        string         `json:"title"`
	Note         sql.NullString `json:"note"`
	UserID       int64          `json:"user_id"`
	BookmarkedAt time.Time      `json:"bookmarked_at"`
	Tags         sql.NullString `json:"tags"`
}

type Token struct {
	ID         int64  `json:"id"`
	TokenHash  string `json:"token_hash"`
	Name       string `json:"name"`
	ShortToken string `json:"short_token"`
	UserID     int64  `json:"user_id"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
