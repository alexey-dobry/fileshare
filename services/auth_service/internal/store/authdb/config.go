package authdb

import (
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb/pg"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb/rd"
)

type Config struct {
	PgConfig    pg.Config `yaml:"pg_config"`
	RedisConfig rd.Config `yaml:"rd_config"`
}
