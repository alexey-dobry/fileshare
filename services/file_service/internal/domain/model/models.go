package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	UUID       string
	Name       string
	MimeType   string
	Size       int64
	UploaderID string

	CourseID string
	GroupID  string

	StorageKey string
	CreatedAt  time.Time
}

type StorageObjectInfo struct {
	Size         int64
	ContentType  string
	LastModified time.Time
}
