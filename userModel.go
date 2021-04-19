package goAuth

import "gorm.io/gorm"

type User struct {
	Username string `gorm:"unique;not null;"`
	Password string `gorm:"not null;"`

	FirstName string
	LasttName string

	Group []*Group `gorm:"many2many:group_users"`

	/*
		.
		.
		.
		Other fields you wants to declare
		to the model
		.
		.
		.
	*/

	gorm.Model
}
