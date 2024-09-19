package handlers

import (
	"github.com/naufalatha/go-boilerplate/config"
	"github.com/naufalatha/go-boilerplate/database"
	"github.com/rs/zerolog"
)

type Handler struct {
}

func InitHandlers(config *config.Configuration, db *database.Database, logger *zerolog.Logger) *Handler {
	return &Handler{}
}
