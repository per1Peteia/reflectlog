package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/per1Peteia/rfl/internal/config"
	"github.com/per1Peteia/rfl/internal/database"
	rflserver "github.com/per1Peteia/rfl/internal/rflserver"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB connection string is not set.")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open database.")
	}

	dbQueries := database.New(db)

	cfg := &config.Cfg{
		DbQueries: dbQueries,
		HttpPort:  "8080",
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: mux,
	}

	mux.HandleFunc("POST /api/clips", rflserver.CreateChirpHandler(cfg))

	log.Printf("Serving on port %s", server.Addr)
	log.Fatal(server.ListenAndServe())

}
