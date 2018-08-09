package models

import (
	"errors"
	"fmt"
	"time"
)

// Address contains general Address
type Address struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Address1  string    `valid:"required" schema:"address1"`
	Address2  string    `valid:"required" schema:"address2"`
	City      string    `valid:"required" schema:"city"`
	State     string    `valid:"required" schema:"state"`
	Zip       string    `valid:"required" schema:"zip"`
	County    string    `valid:"required" schema:"county"`
	Country   string    `valid:"required" schema:"country"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Address) SingleLine() string {
	return fmt.Sprintf("%s %s, %s, %s %s, %s, %s",
		m.Address1,
		m.Address2,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Country,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Address) MultiLine() string {
	return fmt.Sprintf("%s %s\n%s, %s %s\n%s\n%s",
		m.Address1,
		m.Address2,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Country,
	)
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Address) HTMLView() string {
	return "<div id=\"AddressHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Address) HTMLForm() string {
	return "<div id=\"AddressHTMLForm\">{Form Content}</div>"
}

//IsValid returns error if address is not complete
func (m *Address) IsValid() error {
	if m.Address1 == "" || m.City == "" || m.Zip == "" || m.State == "" {
		return errors.New("Please fill in all required fields")
	}
	return nil
}
