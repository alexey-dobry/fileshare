package minio

import (
	"time"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"
	"github.com/minio/minio-go/v7"
)

const maxRetries = 10
const delay = 2 * time.Second

type Repository struct {
	db     *minio.Client
	bucket string
	logger logger.Logger
}

func New(db *minio.Client, logger logger.Logger, bucket string) store.FileRepository {
	return &Repository{
		db:     db,
		bucket: bucket,
		logger: logger,
	}
}
