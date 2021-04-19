package model

import "gorm.io/gorm"

type Group struct {
	Name string `gorm:"not null;unique"`

	Desc string

	Policy []Policy

	User []*User `gorm:"many2many:group_users"`

	gorm.Model
}
