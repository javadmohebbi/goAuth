package goAuth

import "gorm.io/gorm"

type Group struct {
	Name string `gorm:"type:nvarchar(100);not null;unique"`

	Desc string `gorm:"type:nvarchar(100);"`

	Policy []Policy

	User []*User `gorm:"many2many:group_users"`

	gorm.Model
}
