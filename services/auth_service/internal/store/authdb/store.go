package authdb

import (
	"fmt"
	"time"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/model"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb/pg"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb/rd"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 10
const delay = 2 * time.Second

type authStore struct {
	userDB      *gorm.DB
	blacklistDB *redis.Client
	user        store.UserRepository
	blacklist   store.TokenBlacklistRepository
}

func New(logger logger.Logger, cfg Config) (store.Store, error) {

	pgDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.PgConfig.Host,
		cfg.PgConfig.User,
		cfg.PgConfig.Password,
		cfg.PgConfig.DatabaseName,
		cfg.PgConfig.Port,
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

	redisDSN := fmt.Sprintf("redis://%s:%s@%s:%s/%d",
		cfg.RedisConfig.User,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Port,
		cfg.RedisConfig.DatabaseName,
	)

	var redisDB *redis.Client
	opt, err := redis.ParseURL(redisDSN)
	if err != nil {
		return nil, err
	}
	for range maxRetries {
		redisDB = redis.NewClient(opt)
		if err == nil {
			break
		}

		time.Sleep(delay)
	}
	if err != nil {
		return nil, err
	}

	err = pgDB.AutoMigrate(model.Credentials{})
	if err != nil {
		return nil, err
	}

	return &authStore{
		userDB:      pgDB,
		blacklistDB: redisDB,
		user:        pg.New(pgDB, logger),
		blacklist:   rd.New(redisDB, logger),
	}, nil
}

func (as *authStore) User() store.UserRepository {
	return as.user
}

func (as *authStore) Blacklist() store.TokenBlacklistRepository {
	return as.blacklist
}

func (as *authStore) Close() error {
	sqlDB, _ := as.userDB.DB()
	err := sqlDB.Close()
	if err != nil {
		return err
	}

	err = as.blacklistDB.Close()
	if err != nil {
		return err
	}

	return nil
}
