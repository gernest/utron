package models

import (
	"errors"
	"fmt"
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

func (m *PersonName) IsValid() error {
	if m.First == "" || m.Last == "" {
		return errors.New("first or last name can't be empty")
	}

	return nil
}
