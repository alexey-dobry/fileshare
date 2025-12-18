package store

import (
	"io"

	"github.com/alexey-dobry/fileshare/services/file_service/internal/domain/model"
)

type FileRepository interface {
	Put(key string, reader io.Reader, size int64, contentType string) error
	Get(key string) (io.ReadCloser, error)
	Delete(key string) error
	Stat(key string) (*model.StorageObjectInfo, error)
}

type MetaRepository interface {
	Create(file *model.File) error
	GetByID(id string) (*model.File, error)
	Delete(id string) error

	ListByCourse(courseID string) ([]*model.File, error)
	ListByGroup(groupID string) ([]*model.File, error)
}
