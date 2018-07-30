package models

import (
	"errors"
	"fmt"
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
	EmailID        int       `schema:"email_id"`
	Email          Email     `gorm:"foreignkey:EmailID"`
	VerifyPass     string    `gorm:"-" schema:"verifypass"`
	CompanyID      int       `schema:"company_id"`
	Company        Company   `gorm:"foreignkey:CompanyID"`
	PersonID       int       `schema:"person_id"`
	Person         Person    `gorm:"foreignkey:PersonID"`
	HashedPassword string
}

// SingleLine returns a formatted single line text representing the Model
// {Username}: {Email.Address} [{ID},{CID},{PID}]}
func (m *Account) SingleLine() string {
	return fmt.Sprintf("%s: %s [%d, %d, %d]",
		m.Username,
		m.Email.SingleLine(),
		m.ID,
		m.CompanyID,
		m.PersonID,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
// {Username}: {Person.SingleLine()}
// {Email.Address}
// {Company.SingleLine()}
func (m *Account) MultiLine() string {
	return fmt.Sprintf("%s: %s\n%s\n%s\n",
		m.Username,
		m.Person.PersonName.SingleLine(),
		m.Email.Address,
		m.Company.SingleLine(),
	)
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Account) HTMLView() string {
	return "<div id=\"AccountHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Account) HTMLForm() string {
	return "<div id=\"AccountHTMLForm\">{Form Content}</div>"
}

// Validate is used to verifiy password hash match
func (m *Account) Validate() error {
	_, err := govalidator.ValidateStruct(m)
	if err != nil {
		return err
	}
	if m.Password != m.VerifyPass {
		return errors.New("Model.Account: Password missmatch")
	}
	return err
}
