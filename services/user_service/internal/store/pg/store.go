package pg

import (
	"database/sql"
	"fmt"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/course"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/group"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/user"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `validate:"required" mapstructure:"port"`
	User     string `validate:"required" mapstructure:"user"`
	Password string `validate:"required" mapstructure:"password"`
	DB       string `validate:"required" mapstructure:"db"`
}

type pgStore struct {
	db *sql.DB

	user   store.UserRepository
	group  store.GroupRepository
	course store.CourseRepository
}

func New(logger logger.Logger, cfg Config) (store.Store, error) {
	logger = logger.WithFields("layer", "pgstore")

	// pgs connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	// opening sql connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// try to ping db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("pgstore was connected")

	return &pgStore{
		db:     db,
		user:   user.New(logger, db),
		group:  group.New(logger, db),
		course: course.New(logger, db),
	}, nil
}

func (s *pgStore) Close() error {
	return s.db.Close()
}

func (s *pgStore) Course() store.CourseRepository {
	return s.course
}

func (s *pgStore) Group() store.GroupRepository {
	return s.group
}

func (s *pgStore) User() store.UserRepository {
	return s.user
}
