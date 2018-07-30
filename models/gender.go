package models

import (
	"fmt"
	"time"
)

// Gender aims to be LGBT+ compliant and is primarly used for referencing the 'Person'
// in the webapp and templating system
type Gender struct {
	ID         int       `schema:"id"`
	CreatedAt  time.Time `schema:"created"`
	UpdatedAt  time.Time `schema:"updated"`
	ClaimedSex string    `schema:"claimed_sex"` // what they claim -> male, female, gay, lesbian, transgender, etc
	BiosSex    byte      `schema:"legal_sex"`   //What is on birth certificate / under the hood? 0=Unknown, 1=Male, 2=Female
}

// SingleLine returns a formatted single line text representing the Model
func (m *Gender) SingleLine() string {
	return fmt.Sprintf("%s (%s)", m.ClaimedSex, m.BioSexToString(m.BiosSex))
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Gender) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Gender) HTMLView() string {
	return "<div id=\"GenderHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Gender) HTMLForm() string {
	return "<div id=\"GenderHTMLForm\">{Form Content}</div>"
}

// BioSexToString translates the byte value to human readable friendly string
func (m *Gender) BioSexToString(gender byte) string {
	if gender == 1 {
		return "Male"
	} else if gender == 2 {
		return "Female"
	} else {
		return "Unknown"
	}
}
