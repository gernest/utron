package models

import (
	"fmt"
	"time"
)

// Note
type Note struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	PersonID  int       `schema:"person_id"`
	Person    Person    `gorm:"foreignkey:PersonID"`
	Body      string    `schema:"body"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Note) SingleLine() string {
	return fmt.Sprintf("%s, (%s)", m.Body, m.Person.PrimaryName.SingleLine())
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Note) MultiLine() string {
	return fmt.Sprintf("%s:\n%s", m.Person.PrimaryName.SingleLine(), m.Body)
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Note) HTMLView() string {
	return "<div id=\"NoteHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Note) HTMLForm() string {
	return "<div id=\"NoteHTMLForm\">{Form Content}</div>"
}
