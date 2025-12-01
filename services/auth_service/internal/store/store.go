package store

type Store interface {
	User() UserRepository

	Blacklist() TokenBlacklistRepository

	Close() error
}
