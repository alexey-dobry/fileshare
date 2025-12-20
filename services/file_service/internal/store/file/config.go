package file

import (
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/minio"
	"github.com/alexey-dobry/fileshare/services/file_service/internal/store/file/pg"
)

type Config struct {
	PgConfig    pg.Config    `yaml:"pg_config"`
	MinioConfig minio.Config `yaml:"minio_config"`
}
