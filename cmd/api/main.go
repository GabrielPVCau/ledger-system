package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httpHandler "github.com/gabrielcau/ledger-system/internal/handler/http"
	"github.com/gabrielcau/ledger-system/internal/repository/postgres"
	"github.com/gabrielcau/ledger-system/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// If env vars are empty (running locally without docker envs set manually?), fallback or just let it fail.
	// For local dev outside docker, usually we'd have defaults, but let's stick to the container expectation or set defaults if needed.
	if dbHost == "" {
		// Fallback for local run if needed, but docker-compose handles this.
		log.Println("Env vars missing, using local defaults...")
		connStr = "host=localhost port=5432 user=ledger_user password=ledger_pass dbname=ledger_db sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println("Waiting for DB...")
		// In a real app we might retry, but Docker healthcheck helps.
		log.Fatal(err)
	}

	accountRepo := postgres.NewAccountRepository(db)
	ledgerService := service.NewLedgerService(accountRepo)
	transferHandler := httpHandler.NewTransferHandler(ledgerService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/transfer", transferHandler.MakeTransfer)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
