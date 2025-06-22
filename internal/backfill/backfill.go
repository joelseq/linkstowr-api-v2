package backfill

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Link struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Note         string `json:"note"`
	User         string `json:"user"`
	BookmarkedAt string `json:"bookmarked_at"`
}

type Token struct {
	TokenHash  string `json:"token_hash"`
	Name       string `json:"name"`
	ShortToken string `json:"short_token"`
	User       string `json:"user"`
}

func Run(ctx context.Context, db *sql.DB) error {
	users, err := getUsers()
	if err != nil {
		return err
	}

	links, err := getLinks()
	if err != nil {
		return err
	}

	tokens, err := getTokens()
	if err != nil {
		return err
	}

	mapping := buildMapping(users, links, tokens)

	for _, userData := range mapping {
		// Insert user data into the database
		row := db.QueryRowContext(ctx, "INSERT INTO users (username, password) VALUES (?, ?) RETURNING id", userData.Username, userData.Password)
		var userID int64
		err = row.Scan(&userID)
		if err != nil {
			return err
		}

		for _, link := range userData.Links {
			_, err := db.Exec("INSERT INTO links (url, title, note, user_id, bookmarked_at) VALUES (?, ?, ?, ?, ?)",
				link.URL, link.Title, link.Note, userID, link.BookmarkedAt)
			if err != nil {
				return err
			}
		}

		for _, token := range userData.Tokens {
			_, err := db.Exec("INSERT INTO tokens (token_hash, name, short_token, user_id) VALUES (?, ?, ?, ?)",
				token.TokenHash, token.Name, token.ShortToken, userID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type UserData struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Links    []Link  `json:"links"`
	Tokens   []Token `json:"tokens"`
}

func buildMapping(users []User, links []Link, tokens []Token) map[string]UserData {
	mapping := make(map[string]UserData)

	for _, user := range users {
		mapping[user.ID] = UserData{
			Username: user.Username,
			Password: user.Password,
			Links:    []Link{},
			Tokens:   []Token{},
		}
	}

	for _, link := range links {
		if userData, exists := mapping[link.User]; exists {
			userData.Links = append(userData.Links, link)
			mapping[link.User] = userData
		}
	}

	for _, token := range tokens {
		if userData, exists := mapping[token.User]; exists {
			userData.Tokens = append(userData.Tokens, token)
			mapping[token.User] = userData
		}
	}

	return mapping
}

func getUsers() ([]User, error) {
	resp, err := makeRequest("SELECT * FROM user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData []struct {
		Result []User `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData[0].Result, nil
}

func getLinks() ([]Link, error) {
	resp, err := makeRequest("SELECT * FROM link")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData []struct {
		Result []Link `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData[0].Result, nil
}

func getTokens() ([]Token, error) {
	resp, err := makeRequest("SELECT * FROM token")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData []struct {
		Result []Token `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData[0].Result, nil
}

func makeRequest(query string) (*http.Response, error) {
	reqBody := bytes.NewBufferString(query)
	req, err := http.NewRequest(http.MethodPost, "https://linkstowrdb.fly.dev/sql", reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("NS", os.Getenv("SURREALDB_NS"))
	req.Header.Set("DB", os.Getenv("SURREALDB_DB"))
	req.SetBasicAuth(os.Getenv("SURREALDB_USER"), os.Getenv("SURREALDB_PASSWORD"))

	return http.DefaultClient.Do(req)
}
