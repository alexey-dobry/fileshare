package models

import "time"

type Group struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
