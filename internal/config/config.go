package config

import "github.com/per1Peteia/rfl/internal/database"

type Cfg struct {
	DbQueries *database.Queries
	HttpPort  string
}
