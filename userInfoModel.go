package goAuth

import "gorm.io/gorm"

// dynamic fields for user model
type UserInfoField struct {
	FieldName string `gorm:"type:nvarchar(100);unique:not null;"`

	FieldDesc string `gorm:"type:nvarchar(500);default:'-'"`

	// Field Type
	// 1 = number
	// 2 = bool
	// 3 = string
	// 4 = date/time
	// 5 = date
	// 6 = time
	FieldType uint

	NotNull bool

	Unique bool

	Regex string `gorm:"type:ntext;default:'*'"`

	UserInfoFieldDetail []UserInfoFieldDetail

	gorm.Model
}

// dynamic fields for user info
type UserInfoFieldDetail struct {
	Value string `gorm:"type:ntext;"`

	UserInfoFieldID uint

	UserID uint

	gorm.Model
}
