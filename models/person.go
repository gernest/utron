package models

import (
	"fmt"
	"time"
)

//Person contains all data pertaining to a individual person
type Person struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Dob       time.Time `schema:"dob"`

	//NameID is the UID of the Person's name as found in the person_name table
	NameID     int        `schema:"name_id"`
	PersonName PersonName `gorm:"foreignkey:NameID"`

	//EmailID is the UID of the Person's email as found in the email table
	EmailID int   `schema:"email_id"`
	Email   Email `gorm:"foreignkey:EmailID"`

	//TypeID is the UID of the Person's Type as found in the person_type table
	TypeID     int        `schema:"type_id"`
	PersonType PersonType `gorm:"foreignkey:TypeID"`

	//PhoneID is the UID of the Person's Phone info as found in the phone table
	PhoneID int   `schema:"phone_id"`
	Phone   Phone `gorm:"foreignkey:PhoneID"`
}

// SingleLine returns a formatted single line text representing a Person Model
func (m *Person) SingleLine() string {
	pn := m.PersonName.SingleLine()
	return fmt.Sprintf("%s [%s], %s, %s",
		pn,
		m.PersonType.SingleLine(),
		m.Email.Address,
		m.Phone.Number,
	)
}

// MultiLine returns a formatted multi-line text representing a Person Model
func (m *Person) MultiLine() string {
	pn := m.PersonName.SingleLine()
	return fmt.Sprintf("%s\n%s\n%s\n%s\n",
		pn,
		m.PersonType.Name,
		m.Email.Address,
		m.Phone.Number,
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
