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
	Number    int       `valid:"required" schema:"number"`
	Street    string    `valid:"required" schema:"street"`
	City      string    `valid:"required" schema:"city"`
	State     string    `valid:"required" schema:"state"`
	Zip       string    `valid:"required" schema:"zip"`
	County    string    `valid:"required" schema:"county"`
	Country   string    `valid:"required" schema:"country"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Address) SingleLine() string {
	return fmt.Sprintf("%d %s, %s, %s %s, %s, %s",
		m.Number,
		m.Street,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Country,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Address) MultiLine() string {
	return fmt.Sprintf("%d %s\n%s, %s %s\n%s\n%s",
		m.Number,
		m.Street,
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

//TODO Rewrite
func (m *Address) IsValid() error {
	if m.Number <= 0 || m.Street == "" || m.City == "" || m.Zip == "" || m.County == "" || m.Country == "" || m.State == "" {
		return errors.New("one of field is empty or number is negative")
	}

	return nil
}