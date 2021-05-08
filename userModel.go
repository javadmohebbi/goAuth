package goAuth

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"type:nvarchar(100);unique;not null;"`
	Password string `gorm:"not null;"`

	// first name & last name
	// can be provided through fields
	FirstName string `gorm:"type:nvarchar(100);"`
	LastName  string `gorm:"type:nvarchar(100);"`

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

	UserInfoFieldDetail []UserInfoFieldDetail

	MetaData map[string]interface{} `gorm:"-" json:"-"`

	// this ID will
	_id uint `gorm:"-"`

	gorm.Model
}

// This will be called
// after gorm save function
func (u *User) BeforeSave(db *gorm.DB) (err error) {
	return u.beforeSaveAndCreate(db)
}

// This will be called
// after gorm create function
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	return u.beforeSaveAndCreate(db)
}

// This will be called
// after gorm create function
func (u *User) AfterCreate(db *gorm.DB) (err error) {
	for k, v := range u.MetaData {
		var f UserInfoField
		db.Where("field_name = ?", k).First(&f)
		if f.ID != 0 {

			value := fmt.Sprintf("%s", v)

			fd := UserInfoFieldDetail{
				UserInfoFieldID: f.ID,
				UserID:          u.ID,
				Value:           value,
			}

			db.Create(&fd)

			u.UserInfoFieldDetail = append(u.UserInfoFieldDetail, fd)

		}
	}

	return nil
}

// run after find
func (u *User) AfterFind(db *gorm.DB) (err error) {
	// filter password fileds
	u.Password = "__FILTERED__"
	u.UserInfoFieldDetail = []UserInfoFieldDetail{}

	var fds []UserInfoFieldDetail
	m := make(map[string]interface{})

	db.Where("user_id = ?", u.ID).Find(&fds)
	for _, fd := range fds {
		var f UserInfoField
		db.Where("id = ?", fd.UserInfoFieldID).First(&f)
		if f.ID != 0 {
			m[f.FieldName] = fd.Value
			u.UserInfoFieldDetail = append(u.UserInfoFieldDetail, fd)
		}
	}

	return nil
}

// this whill be called before save & update
func (u *User) beforeSaveAndCreate(db *gorm.DB) (err error) {

	// hash the password
	b, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	hashed := string(b)
	u.Password = hashed

	// validate & return
	return u.metaValidator(db)
}

// validate user metda data validation
func (u *User) metaValidator(db *gorm.DB) (err error) {

	// username validation
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]+((_|-|\.)?[a-zA-Z0-9])*$`, u.Username)
	if !match {
		return errors.New(fmt.Sprintf("Invalid username '%s'", u.Username))
	}

	// password validation
	if psErr := u.passwordValidator(); psErr != nil {
		return psErr
	}

	for k, v := range u.MetaData {
		var f UserInfoField
		db.Where("field_name = ?", k).First(&f)
		if f.ID != 0 {
			value := fmt.Sprintf("%s", v)

			// continue in case not null is true and value is empty
			if f.NotNull == true && value == "" {
				continue
			}

			// check if duplicated
			if f.Unique == true {

				// uniq can not have empty
				if value == "" {
					return errors.New(fmt.Sprintf("Uniue field '%s' can not have empty value", k))
				}

				var fd UserInfoFieldDetail
				db.Where("value = ? AND user_info_field_id = ?", value, f.ID).First(&fd)
				if fd.ID != 0 {
					return errors.New(fmt.Sprintf("Value of field '%s' must be unique, but the provided value '%s' is duplicated!", k, value))
				}
			}

			// validate
			re := regexp.MustCompile(f.Regex)
			if !re.MatchString(value) {
				// invalid
				return errors.New(fmt.Sprintf("Invalid value (%v) for field %s", value, f.FieldName))
			}
		} else {
			return errors.New(fmt.Sprintf("Uknown field name '%s'", k))
		}
	}

	return nil
}

// validate password
func (u *User) passwordValidator() (err error) {
	letters := 0
	specials := 0
	numbers := 0

	checkNum, checkUpper, checkSpecial, checkLetter := false, false, false, false

	if len(u.Password) < 8 || len(u.Password) > 20 {
		return errors.New("Password length be between 8 and 20 character")
	}

	for _, c := range u.Password {
		switch {
		case unicode.IsNumber(c):
			checkNum = true
			numbers++
		case unicode.IsUpper(c):
			checkUpper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			checkSpecial = true
			specials++
		case unicode.IsLetter(c) || c == ' ':
			checkLetter = true
			letters++
		default:
			//return false, false, false, false
		}
	}

	if !checkNum {
		return errors.New("Password must contain at least 1 number")
	}

	if !checkUpper {
		return errors.New("Password must contain at least 1 capital letter")
	}

	if !checkLetter {
		return errors.New("Password must contain at least 1 letter")
	}

	if !checkSpecial {
		return errors.New("Password must contain at least 1 special characters")
	}

	return nil

}
