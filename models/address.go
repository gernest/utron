package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	nsmisc "github.com/NlaakStudios/gowaf/utils/misc"
)

//TODO I've deleted valid:required, because of using govalidator for account
// Address contains general Address
type Address struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Address1  string    `schema:"address1"`
	Address2  string    `schema:"address2"`
	City      string    `schema:"city"`
	State     string    `schema:"state"`
	Zip       string    `schema:"zip"`
	County    string    `schema:"county"`
	Country   string    `schema:"country"`
	Friendly  string    `schema:"friendly"`
}

// SingleLine returns a formatted single line text representing the Model
func (m *Address) SingleLine() string {
	return fmt.Sprintf("%s %s, %s, %s %s, %s, %s",
		m.Address1,
		m.Address2,
		m.City,
		m.State,
		m.Zip,
		m.County,
		m.Country,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Address) MultiLine() string {
	return fmt.Sprintf("%s %s\n%s, %s %s\n%s\n%s",
		m.Address1,
		m.Address2,
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

// Sanitize strips all leading and trailing whitespace from strings as well as test normalization all model string properties.
func (m *Address) Sanitize() {
	if m.Country == "" {
		m.Country = "United States"
	}
	m.Address1 = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.Address1)))
	m.Address2 = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.Address2)))
	m.City = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.City)))
	m.State = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.State)))
	m.Zip = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.Zip)))
	m.County = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.County)))
	m.Country = strings.ToTitle(strings.TrimSpace(nsmisc.StripCtlAndExtFromUTF8(m.Country)))
	m.Friendly = strings.TrimSpace(m.SingleLine())
}

//IsValid returns error if model is not complete
func (m *Address) IsValid() error {
	if m.Address1 == "" || m.City == "" || m.Zip == "" || m.State == "" {
		return errors.New("Please fill in all required fields")
	}
	return nil
}