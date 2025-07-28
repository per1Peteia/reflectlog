package main

import "github.com/per1Peteia/rfl/internal/database"

const PORT = "8080"

type Cfg struct {
	dbQueries *database.Queries
	addr      string
}
