package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

//PersonName hold the complete name of a person
type PersonName struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Prefix    string    `schema:"prefix"`  //ie. Mr
	First     string    `schema:"first"`   //William
	Middle    string    `schema:"middle"`  //Blaine
	Last      string    `schema:"last"`    //Doe
	Suffix    string    `schema:"suffix"`  //Sr
	GoesBy    string    `schema:"goes_by"` //Bob
	Friendly  string    `schema:"friendly"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *PersonName) SingleLine() string {
	return fmt.Sprintf("%s %s %s %s %s\n",
		m.Prefix,
		m.First,
		m.Middle,
		m.Last,
		m.Suffix,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *PersonName) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *PersonName) HTMLView() string {
	return "<div id=\"PersonNameHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *PersonName) HTMLForm() string {
	return "<div id=\"PersonNameHTMLForm\">{Form Content}</div>"
}

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *PersonName) Sanitize() {
	m.Prefix = strings.ToTitle(strings.TrimSpace(m.Prefix))
	m.First = strings.ToTitle(strings.TrimSpace(m.First))
	m.Middle = strings.ToTitle(strings.TrimSpace(m.Middle))
	m.Last = strings.ToTitle(strings.TrimSpace(m.Last))
	m.Suffix = strings.ToTitle(strings.TrimSpace(m.Suffix))
	m.GoesBy = strings.ToTitle(strings.TrimSpace(m.GoesBy))
	m.Friendly = strings.ToTitle(strings.TrimSpace(m.SingleLine()))
}

func (m *PersonName) IsValid() error {
	if m.First == "" || m.Last == "" {
		return errors.New("first or last name can't be empty")
	}

	return nil
}
