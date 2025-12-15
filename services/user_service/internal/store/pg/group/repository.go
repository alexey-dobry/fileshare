package group

import (
	"database/sql"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store"
)

type Repository struct {
	db     *sql.DB
	logger logger.Logger
}

func New(logger logger.Logger, db *sql.DB) store.GroupRepository {

	return &Repository{
		db:     db,
		logger: logger.WithFields("layer", "group repository"),
	}
}
