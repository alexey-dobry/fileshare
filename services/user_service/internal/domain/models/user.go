package models

import "time"

type User struct {
	ID        string
	Name      string
	Surname   string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
