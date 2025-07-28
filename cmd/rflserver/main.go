package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/per1Peteia/rfl/internal/database"
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

	c := &Cfg{
		dbQueries,
		":" + PORT,
	}

	m := http.NewServeMux()
	s := &http.Server{
		Addr:    c.addr,
		Handler: m,
	}

	m.HandleFunc("POST /api/clips", CreateClipHandler(c))

	log.Printf("serving on port %s\n", PORT)
	log.Fatal(s.ListenAndServe())

}
