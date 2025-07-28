package main

import "github.com/per1Peteia/rfl/internal/database"

const PORT = "8080"

type Config struct {
	dbQueries *database.Queries
	addr      string
}
