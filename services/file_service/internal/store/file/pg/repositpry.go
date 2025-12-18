package pg

import (
	"time"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"
	"gorm.io/gorm"
)

const maxRetries = 10
const delay = 2 * time.Second

type Repository struct {
	db     *gorm.DB
	logger logger.Logger
}

func New(db *gorm.DB, logger logger.Logger) store.MetaRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
