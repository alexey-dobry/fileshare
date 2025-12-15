package models

import "time"

type Course struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
