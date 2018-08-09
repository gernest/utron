package models

import (
	"errors"
	"fmt"
	"time"
)

const (
	//PhoneTypeUnknown represents a defaul unknown phone type
	PhoneTypeUnknown = byte(0)
	//PhoneTypeMobile represents a Mobile or Cell phone number
	PhoneTypeMobile = byte(1)
	//PhoneTypeHome represents a home phone number
	PhoneTypeHome = byte(2)
	//PhoneTypeBusiness represents a business phone number
	PhoneTypeBusiness = byte(3)
	//PhoneTypeFax represents a Fax phone number
	PhoneTypeFax = byte(4)
)

//PhoneStruct is used to breakdown and store phone numbers
type Phone struct {
	ID          int       `schema:"id"`
	CreatedAt   time.Time `schema:"created"`
	UpdatedAt   time.Time `schema:"updated"`
	CountryCode string    `schema:"code"`
	AreaCode    string    `schema:"area"`
	Number      string    `schema:"number"`
	PhoneType   byte      `schema:"phone_type"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Phone) SingleLine() string {
	return fmt.Sprintf("%s (%s) %s [%s]", m.CountryCode, m.AreaCode, m.Number, m.PhoneTypeToString(m.PhoneType))
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Phone) MultiLine() string {
	return fmt.Sprintf("%s\n%s (%s) %s\n", m.PhoneTypeToString(m.PhoneType), m.CountryCode, m.AreaCode, m.Number)
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Phone) HTMLView() string {
	return "<div id=\"PhoneHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Phone) HTMLForm() string {
	return "<div id=\"PhoneHTMLForm\">{Form Content}</div>"
}

/***[Support Methods]***/

// PhoneTypeToString given a valid PhoneType Byte value will return the string representation
func (m *Phone) PhoneTypeToString(pt byte) string {
	switch {
	case pt == PhoneTypeUnknown:
		return "Unknown"
	case pt == PhoneTypeMobile:
		return "Mobile"
	case pt == PhoneTypeHome:
		return "Home"
	case pt == PhoneTypeBusiness:
		return "Business"
	case pt == PhoneTypeFax:
		return "Fax"
	}
	return ""
}

func (m *Phone) IsValid() error {
	if m.CountryCode == "" || m.Number == "" {
		return errors.New("country code and number can't be empty")
	}

	return nil
}
