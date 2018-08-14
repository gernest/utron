package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//Person contains all data pertaining to a individual person
type Person struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Dob       time.Time `schema:"dob"`

	//GenderID is the UID of the Person's Gender as found in the gender table
	GenderID int    `schema:"gender_id"`
	Gender   Gender `gorm:"foreignkey:GenderID"`

	//NameID is the UID of the Person's name as found in the person_name table
	NameID      int        `schema:"name_id"`
	PrimaryName PersonName `gorm:"foreignkey:NameID"`

	//SpouseNameID is the UID of the Person's name as found in the person_name table
	SpouseNameID int        `schema:"spouse_name_id"`
	SpouseName   PersonName `gorm:"foreignkey:SpouseNameID"`

	//EmailID is the UID of the Person's email as found in the email table
	EmailID int   `schema:"email_id"`
	Email   Email `gorm:"foreignkey:EmailID"`

	//TypeID is the UID of the Person's Type as found in the person_type table
	TypeID     int        `schema:"type_id"`
	PersonType PersonType `gorm:"foreignkey:TypeID"`

	//PhoneID is the UID of the Person's Phone info as found in the phone table
	PhoneID  int    `schema:"phone_id"`
	Phone    Phone  `gorm:"foreignkey:PhoneID"`
	Friendly string `schema:"friendly"`
}

// SingleLine returns a formatted single line text representing a Person Model
func (m *Person) SingleLine() string {
	pn := m.PrimaryName.SingleLine()
	return fmt.Sprintf("%s [%s], %s, %s",
		pn,
		m.PersonType.SingleLine(),
		m.Email.Friendly,
		m.Phone.Number,
	)
}

// MultiLine returns a formatted multi-line text representing a Person Model
func (m *Person) MultiLine() string {
	pn := m.PrimaryName.SingleLine()
	return fmt.Sprintf("%s\n%s\n%s\n%s\n",
		pn,
		m.PersonType.Name,
		m.Email.Friendly,
		m.Phone.Friendly,
	)
}

// HTMLView returns a HTML5 code representing a view of a Person Model
func (m *Person) HTMLView() string {
	return "<div id=\"PersonHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of a Person Model
func (m *Person) HTMLForm() string {
	return "<div id=\"PersonHTMLForm\">{Form Content}</div>"
}

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *Person) Sanitize() {
	m.Friendly = strings.TrimSpace(m.SingleLine())
}

//IsValid returns error if model is not complete
func (m *Person) IsValid() error {
	if m.NameID < 1 || m.EmailID < 1 || m.PhoneID < 1 {
		return errors.New("Please fill in all required fields")
	}
	return nil
}
