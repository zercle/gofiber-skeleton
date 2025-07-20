package entities

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`	
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
