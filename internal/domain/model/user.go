package model

import "time"

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	RoleId    string    `json:"roleId"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
