package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Example contains general Example
type Example struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Number    int       `schema:"number"`                  //123456
	String    string    `valid:"required" schema:"string"` //Test String
	Toggle    bool      `schema:"toggle"`                  //True/False, On/Off, Yes/No
	Float     float64   `schema:"float"`                   //$12.56, 1,123,456.987654321
}

// SingleLine returns a formatted single line text representing the Model
func (m *Example) SingleLine() string {
	return fmt.Sprintf("%d %s, %s, %s",
		m.Number,
		m.String,
		strconv.FormatBool(m.Toggle),
		strconv.FormatFloat(m.Float, 'E', -1, 64),
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Example) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Example) HTMLView() string {
	return "<div id=\"ExampleHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Example) HTMLForm() string {
	return "<div id=\"ExampleHTMLForm\">{Form Content}</div>"
}

//IsValid returns error if model is not complete
func (m *Example) IsValid() error {
	if m.String == "" {
		return errors.New("Please fill in all required fields")
	}
	return nil
}
