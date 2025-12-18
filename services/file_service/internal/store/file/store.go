package file

import (
	"context"
	"fmt"
	"time"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/domain/model"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store"
	mn "github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/minio"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/pg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 10
const delay = 2 * time.Second

type authStore struct {
	metaDB *gorm.DB
	fileDB *minio.Client
	meta   store.MetaRepository
	file   store.FileRepository
}

func New(logger logger.Logger, cfg Config) (store.Store, error) {
	pgDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.pgConfig.Host,
		cfg.pgConfig.User,
		cfg.pgConfig.Password,
		cfg.pgConfig.DatabaseName,
		cfg.pgConfig.Port,
	)

	var pgDB *gorm.DB
	var err error
	for range maxRetries {
		pgDB, err = gorm.Open(postgres.Open(pgDSN))
		if err == nil {
			break
		}

		time.Sleep(delay)
	}
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s:%s", cfg.minioConfig.Host, cfg.minioConfig.Port)
	accessKey := cfg.minioConfig.AccessKey
	secretAccessKey := cfg.minioConfig.SecretKey
	useSSL := false

	minioDB, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioDB.BucketExists(context.Background(), cfg.minioConfig.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := minioDB.MakeBucket(context.Background(), cfg.minioConfig.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	err = pgDB.AutoMigrate(model.File{})
	if err != nil {
		return nil, err
	}

	return &authStore{
		metaDB: pgDB,
		fileDB: minioDB,
		meta:   pg.New(pgDB, logger),
		file:   mn.New(minioDB, logger, cfg.minioConfig.Bucket),
	}, nil
}

func (as *authStore) Meta() store.MetaRepository {
	return as.meta
}

func (as *authStore) File() store.FileRepository {
	return as.file
}

func (as *authStore) Close() error {
	sqlDB, _ := as.metaDB.DB()
	err := sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
