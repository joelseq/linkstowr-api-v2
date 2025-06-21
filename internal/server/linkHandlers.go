package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"linkstowr/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) listLinksHandler(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	links, err := s.repository.ListLinks(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	if len(links) == 0 {
		links = make([]repository.ListLinksRow, 0)
	}

	return c.JSON(http.StatusOK, links)
}

func (s *Server) createLinkHandler(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	var createLinkPayload struct {
		URL   string `json:"url" validate:"required,url"`
		Title string `json:"title" validate:"required"`
		Note  string `json:"note"`
	}

	err = json.NewDecoder(c.Request().Body).Decode(&createLinkPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	v := validator.New()
	err = v.Struct(createLinkPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	row, err := s.repository.CreateLink(c.Request().Context(), repository.CreateLinkParams{
		UserID: userID,
		Url:    createLinkPayload.URL,
		Title:  createLinkPayload.Title,
		Note:   sql.NullString{String: createLinkPayload.Note, Valid: true},
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"result": echo.Map{
			"url":     row.Url,
			"success": true,
		},
	})
}

func (s *Server) clearLinksHandler(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	err = s.repository.ClearLinks(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}
