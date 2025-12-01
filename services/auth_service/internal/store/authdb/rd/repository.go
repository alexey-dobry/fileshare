package rd

import (
	"time"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
	"github.com/redis/go-redis/v9"
)

const maxRetries = 10
const delay = 2 * time.Second

type Repository struct {
	db     *redis.Client
	logger logger.Logger
}

func New(db *redis.Client, logger logger.Logger) store.TokenBlacklistRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
