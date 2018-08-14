package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
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
	Friendly  string    `schema:"friendly"`
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

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *PersonType) Sanitize() {
	m.Name = strings.ToTitle(strings.TrimSpace(m.Name))
	m.Friendly = strings.ToTitle(strings.TrimSpace(m.Friendly))
}

//IsValid returns error if model is not complete
func (m *PersonType) IsValid() error {
	if m.Name == "" {
		return errors.New("name can't be empty")
	}

	return nil
}
