package server

import (
	"net/http"
	"os"

	"github.com/joelseq/sqliteadmin-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prettylogger "github.com/rdbell/echo-pretty-logger"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(prettylogger.Logger)
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// SQLite admin setup
	config := sqliteadmin.Config{
		DB:       s.db.GetDB(),
		Username: os.Getenv("SQLITE_ADMIN_USERNAME"),
		Password: os.Getenv("SQLITE_ADMIN_PASSWORD"),
	}
	admin := sqliteadmin.New(config)

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)
	e.POST("/admin", wrappedHandler(admin.HandlePost))

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func wrappedHandler(handlerFunc func(w http.ResponseWriter, r *http.Request)) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response().Writer
		r := c.Request()

		// Call the handler function and capture any error
		handlerFunc(w, r)
		return nil
	}
}
