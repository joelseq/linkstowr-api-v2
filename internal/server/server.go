package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"linkstowr/internal/database"
	"linkstowr/internal/repository"
)

type Server struct {
	port int

	db database.Service

	repository *repository.Queries
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	if err := db.RunMigrations(); err != nil {
		log.Fatal(err)
	}

	NewServer := &Server{
		port: port,

		db: database.New(),

		repository: repository.New(db.GetDB()),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
