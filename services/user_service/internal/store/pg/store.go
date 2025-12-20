package pg

import (
	"database/sql"
	"fmt"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/course"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/group"
	"github.com/alexey-dobry/fileshare/services/user_service/internal/store/pg/user"

	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `validate:"required" yaml:"port"`
	User     string `validate:"required" yaml:"user"`
	Password string `validate:"required" yaml:"password"`
	DB       string `yaml:"database"`
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DB,
		cfg.Port,
	)

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

	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	logger.Info("migrations are up")

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
