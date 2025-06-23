package server

import (
	"net/http"
	"os"

	"github.com/joelseq/sqliteadmin-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prettylogger "github.com/rdbell/echo-pretty-logger"

	"linkstowr/internal/auth"
	"linkstowr/internal/backfill"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(prettylogger.Logger)
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

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

	// Auth routes
	e.POST("/signup", s.signupHandler)
	e.POST("/signin", s.signinHandler)
	e.GET("/me", s.meHandler)

	// Backfill
	backfillGroup := e.Group("/backfill")
	backfillGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("BACKFILL_USERNAME") && password == os.Getenv("BACKFILL_PASSWORD") {
			return true, nil
		}
		return false, nil
	}))
	backfillGroup.POST("/run", s.backfillHandler)

	// API routes
	api := e.Group("/api")
	api.Use(auth.GetMiddleware(s.repository))

	// Token routes
	api.GET("/tokens", s.listTokensHandler)
	api.POST("/tokens", s.createTokenHandler)
	api.DELETE("/tokens/:id", s.deleteTokenHandler)

	// Link routes
	api.GET("/links", s.listLinksHandler)
	api.POST("/links", s.createLinkHandler)
	api.POST("/links/clear", s.clearLinksHandler)

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

func (s *Server) backfillHandler(c echo.Context) error {
	err := backfill.Run(c.Request().Context(), s.db.GetDB())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Backfill completed successfully"})
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
