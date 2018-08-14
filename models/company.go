package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//Company stores information about the company
type Company struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Name      string    `schema:"name"`
	ContactID int       `schema:"contact_id"`
	Person    Person    `gorm:"foreignkey:ContactID"`
	PhoneID   int       `schema:"phone_id"`
	Phone     Phone     `gorm:"foreignkey:PhoneID"`
	FaxID     int       `schema:"fax_id"`
	Fax       Phone     `gorm:"foreignkey:FaxID"`
	Friendly  string    `schema:"friendly"`
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

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *Company) Sanitize() {
	m.Name = strings.ToTitle(strings.TrimSpace(m.Name))
	m.Friendly = strings.ToTitle(strings.TrimSpace(m.SingleLine()))
}

//IsValid returns error if model is not complete
func (m *Company) IsValid() error {
	if m.Name == "" || m.ContactID < 1 || m.PhoneID < 1 {
		return errors.New("Please fill in all required fields")
	}
	return nil
}
