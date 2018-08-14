package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
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
	CountryCode string    `schema:"country_code"`
	AreaCode    string    `schema:"area_code"`
	Number      string    `schema:"number"`
	Extension   string    `schema:"extention"`
	PhoneType   byte      `schema:"phone_type"`
	Friendly    string    `schema:"friendly"`
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

// Parse takes a phone number as a string and parses it into the model
func (m *Phone) Parse(p string) {
	parsedPhone, err := phonenumbers.Parse(p, "US")
	if err == nil {
		m.PhoneType = PhoneTypeUnknown
		m.CountryCode = strconv.FormatInt(int64(parsedPhone.GetCountryCode()), 10)
		m.Number = strconv.FormatInt(int64(parsedPhone.GetNationalNumber()), 10)
		m.Extension = parsedPhone.GetExtension()
		if len(m.Number) == 10 {
			m.AreaCode = m.Number[0:3]
			m.Number = fmt.Sprintf("%s-%s", m.Number[3:6], m.Number[6:10])
		}
	}
}

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *Phone) Sanitize() {
	m.CountryCode = strings.ToTitle(strings.TrimSpace(m.CountryCode))
	m.AreaCode = strings.ToTitle(strings.TrimSpace(m.AreaCode))
	m.Number = strings.ToTitle(strings.TrimSpace(m.Number))
	m.Extension = strings.ToTitle(strings.TrimSpace(m.Extension))
	m.Friendly = strings.TrimSpace(m.SingleLine())
}

//IsValid returns error if model is not complete
func (m *Phone) IsValid() error {
	if m.CountryCode == "" || m.Number == "" {
		return errors.New("country code and number can't be empty")
	}

	return nil
}
