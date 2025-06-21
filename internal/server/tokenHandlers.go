package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"linkstowr/internal/auth"
	"linkstowr/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) listTokensHandler(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	tokens, err := s.repository.ListTokens(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		tokens = make([]repository.ListTokensRow, 0)
	}

	return c.JSON(http.StatusOK, tokens)
}

func (s *Server) createTokenHandler(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	var createTokenPayload struct {
		Name string `json:"name" validate:"required"`
	}

	err = json.NewDecoder(c.Request().Body).Decode(&createTokenPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	v := validator.New()
	err = v.Struct(createTokenPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	key, err := auth.NewPrefixedAPIKey()
	if err != nil {
		return err
	}

	err = s.repository.CreateToken(c.Request().Context(), repository.CreateTokenParams{
		TokenHash:  key.LongTokenHash(),
		ShortToken: key.ShortToken(),
		UserID:     userID,
		Name:       createTokenPayload.Name,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": key.Token(),
	})
}

func (s *Server) deleteTokenHandler(c echo.Context) error {
	tokenID := c.Param("id")
	if tokenID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Token ID is required")
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	// Convert tokenID to int64
	id, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token ID")
	}

	err = s.repository.DeleteToken(c.Request().Context(), repository.DeleteTokenParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func getUserIDFromContext(c echo.Context) (int64, error) {
	userID := c.Get("userID").(string)

	return strconv.ParseInt(userID, 10, 64)
}
