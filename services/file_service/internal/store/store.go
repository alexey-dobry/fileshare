package store

type Store interface {
	File() FileRepository

	Meta() MetaRepository

	Close() error
}
