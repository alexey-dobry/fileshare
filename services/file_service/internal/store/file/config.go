package file

import (
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/minio"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/pg"
)

type Config struct {
	pgConfig    pg.Config    `yaml:"pg_config"`
	minioConfig minio.Config `yaml:"minio_config"`
}
