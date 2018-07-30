package models

import "time"

// Gender aims to be LGBT+ compliant and is primarly used for referencing the 'Person'
// in the webapp and templating system
type Gender struct {
	ID         int       `schema:"id"`
	CreatedAt  time.Time `schema:"created"`
	UpdatedAt  time.Time `schema:"updated"`
	ClaimedSex string    `schema:"claimed_sex"` // what they claim -> male, female, gay, lesbian, transgender, etc
	BiosSex    byte      `schema:"legal_sex"`   //What is on birth certificate / under the hood? 0=Unknown, 1=Male, 2=Female
}

// BioSexToString translates the byte value to human readable friendly string
func (g *Gender) BioSexToString(gender byte) string {
	if gender == 1 {
		return "Male"
	} else if gender == 2 {
		return "Female"
	} else {
		return "Unknown"
	}
}