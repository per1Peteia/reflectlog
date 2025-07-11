package server

import (
	"net/http"

	"github.com/per1Peteia/rfl/internal/config"
)

func CreateChirpHandler(cfg *config.Cfg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
