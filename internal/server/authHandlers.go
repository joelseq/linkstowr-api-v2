package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"linkstowr/internal/auth"
	"linkstowr/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) signupHandler(c echo.Context) error {
	var signupPayload struct {
		Username        string `json:"username" validate:"required"`
		Password        string `json:"password" validate:"required"`
		PasswordConfirm string `json:"password_confirm" validate:"required"`
	}

	err := json.NewDecoder(c.Request().Body).Decode(&signupPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	v := validator.New()
	err = v.Struct(signupPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	if signupPayload.Password != signupPayload.PasswordConfirm {
		return echo.NewHTTPError(http.StatusBadRequest, "Passwords do not match")
	}

	hashedPassword, err := auth.HashPassword(signupPayload.Password, auth.DefaultParams)
	if err != nil {
		return err
	}

	row, err := s.repository.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Username: signupPayload.Username,
		Password: hashedPassword,
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return echo.NewHTTPError(http.StatusConflict, "Username already exists")
		}

		return err
	}

	token, err := auth.GenerateJWT(row.ID, row.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"id":       row.ID,
		"username": row.Username,
		"token":    token,
	})
}

func (s *Server) signinHandler(c echo.Context) error {
	var signinPayload struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := json.NewDecoder(c.Request().Body).Decode(&signinPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	v := validator.New()
	err = v.Struct(signinPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	row, err := s.repository.GetUser(c.Request().Context(), signinPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
		}

		return err
	}

	var ok bool

	if ok, err = auth.ComparePasswordAndHash(signinPayload.Password, row.Password); err != nil {
		return err
	} else if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	// Valid credentials, generate JWT
	token, err := auth.GenerateJWT(row.ID, row.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":       row.ID,
		"username": row.Username,
		"token":    token,
	})
}

func (s *Server) meHandler(c echo.Context) error {
	// Get the Authorization header from the request and decode the jwt
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := auth.DecodeJWT(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":       claims.Subject,
		"username": claims.Username,
	})
}
