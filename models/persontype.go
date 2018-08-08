package models

import (
	"fmt"
	"time"
	"errors"
)

// PersonType provides a list of avilable "title" such as:
// 0	Unknown
// 1	Adjuster
// 2	Property Owner
// 3	Attorney
// 4	Paralegal
// 5	Contractor
type PersonType struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Name      string    `schema:"name"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *PersonType) SingleLine() string {
	return fmt.Sprintf("%s",
		m.Name,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *PersonType) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *PersonType) HTMLView() string {
	return "<div id=\"PersonTypeHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *PersonType) HTMLForm() string {
	return "<div id=\"PersonTypeHTMLForm\">{Form Content}</div>"
}

func (m *PersonType) IsValid() error {
	if m.Name == "" {
		return errors.New("name can't be empty")
	}

	return nil
}