package model

import "time"

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
}
