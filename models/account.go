package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
)

// Account is used to represent a user for authentication
type Account struct {
	ID             int       `schema:"id"`
	CreatedAt      time.Time `schema:"created"`
	UpdatedAt      time.Time `schema:"updated"`
	Username       string    `valid:"required,length(6|16)" schema:"username"`
	Password       string    `gorm:"-" valid:"required,length(6|24)" schema:"password"`
	Email          string    `valid:"required,email" schema:"email"`
	VerifyPass     string    `gorm:"-" schema:"verifypass"`
	HashedPassword string
	CompanyID      int `schema:"company_id"`
	PersonID       int `schema:"person_id"`
}

// Validate is used to verifiy password hash match
func (u *Account) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	if u.Password != u.VerifyPass {
		return errors.New("Model.Account: Password missmatch")
	}
	return err
}
