package models

import "time"

//TODO: Remove this from public repo

//Email contains a breakdown of a email
//TODO: Update to use/integrate "net/mail" and Address
type Email struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Address   string    `schema:"address"`  //bob1234@gmail.com
	Username  string    `schema:"username"` //bob1234
	Domain    string    `schema:"domain"`   //gmail.com
}

// SingleLine returns a formatted single line text representing the Model
func (m *Email) SingleLine() string {
	return m.Address
}

// MultiLine returns a formatted multi-line text representing the Model
func (m *Email) MultiLine() string {
	return m.SingleLine()
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Email) HTMLView() string {
	return "<div id=\"EmailHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Email) HTMLForm() string {
	return "<div id=\"EmailHTMLForm\">{Form Content}</div>"
}
