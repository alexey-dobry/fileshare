package store

type Store interface {
	// User returns user repository
	User() UserRepository
	// Manager returns manager repository
	Group() GroupRepository

	//Course returns course repository
	Course() CourseRepository

	// Close closes store
	Close() error
}
