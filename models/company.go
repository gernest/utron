package models

import (
	"fmt"
	"time"
)

//CompanyStruct stores information about the company behind the coin/blockchain
type Company struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Name      string    `schema:"name"`
	ContactID Person    `schema:"contact_id"`
	Person    Person    `gorm:"foreignkey:ContactID"`
	PhoneID   int       `schema:"phone_id"`
	Phone     Phone     `gorm:"foreignkey:PhoneID"`
	FaxID     int       `schema:"fax_id"`
	Fax       Phone     `gorm:"foreignkey:FaxID"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Company) SingleLine() string {
	return fmt.Sprintf("%s", m.Name)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Company) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Company) HTMLView() string {
	return "<div id=\"CompanyHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Company) HTMLForm() string {
	return "<div id=\"CompanyHTMLForm\">{Form Content}</div>"
}
