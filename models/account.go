package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

const (
	/***************************/
	/***[User Account States]***/
	/***************************/

	// UserStateVerifyEmailSend -> Need to send out email verification
	UserStateVerifyEmailSend = 1
	// UserStateVerifyEmailSent -> Email verification has been sent
	UserStateVerifyEmailSent = 2
	// UserStateVerifyEmailDone -> User clicked on link in email, Verified, Account Active
	UserStateVerifyEmailDone = 3
	// UserStateBanned -> Account has been banned. No access granted other than guest
	UserStateBanned = 4
	// UserStateIdle -> User Signed in but no activity for 5 min (300 sec)
	UserStateIdle = 5
	// UserStateSignedIn -> User is Signed In
	UserStateSignedIn = 6
	// UserStateSignedOut -> User is Signed Out
	UserStateSignedOut = 7

	/**************************/
	/***[User Access Rights]***/
	/**************************/

	// UserAccessGuest -> default, no account, Guest access
	UserAccessGuest = 0
	// UserAccessMember -> Active account with access to Member content
	UserAccessMember = 1
	// UserAccessEmployee -> Active account with access to Member and Employee content
	UserAccessEmployee = 2
	// UserAccessAdmin -> Active account with access to Member, Employee and Admin content
	UserAccessAdmin = 3
)

// Account is used to represent a user for authentication
type Account struct {
	ID             int       `schema:"id"`
	CreatedAt      time.Time `schema:"created"`
	UpdatedAt      time.Time `schema:"updated"`
	Username       string    `valid:"required,length(6|16)" schema:"username"`
	EmailAccount   string    `valid:"required,length(6|16), email" schema:"email"`
	Password       string    `gorm:"-" valid:"required,length(6|24)" schema:"password"`
	VerifyPassword string    `gorm:"-" valid:"required,length(6|24)" schema:"verify_password"`
	HashedPassword string    `schema:"hashed_password"`
	State          byte      `schema:"state"`
	Access         byte      `schema:"access"`

	EmailID    int     `schema:"email_id"`
	Email      Email   `gorm:"foreignkey:EmailID"`
	VerifyPass string  `gorm:"-" schema:"verifypass"`
	CompanyID  int     `schema:"company_id"`
	Company    Company `gorm:"foreignkey:CompanyID"`
	PersonID   int     `schema:"person_id"`
	Person     Person  `gorm:"foreignkey:PersonID"`
}

// SingleLine returns a formatted single line text representing the Model
// {Username}: {Email.Address} [{ID},{CID},{PID}]}
func (m *Account) SingleLine() string {
	return fmt.Sprintf("%s: %s [%d, %d, %d]",
		m.Username,
		m.EmailAccount,
		m.ID,
		m.CompanyID,
		m.PersonID,
	)
}

// MultiLine returns a formatted multi-line text representing the Model
// {Username}: {Person.SingleLine()}
// {Email.Address}
// {Company.SingleLine()}
func (m *Account) MultiLine() string {
	return fmt.Sprintf("%s: %s\n%s\n%s\n",
		m.Username,
		m.Person.PrimaryName.SingleLine(),
		m.EmailAccount,
		m.Company.SingleLine(),
	)
}

// HTMLView returns a HTML5 code representing a view of the Model
func (m *Account) HTMLView() string {
	return "<div id=\"AccountHTMLView\">{View Content}</div>"
}

// HTMLForm returns a HTML5 code representing a form of the Model
func (m *Account) HTMLForm() string {
	return "<div id=\"AccountHTMLForm\">{Form Content}</div>"
}

// Validate is used to verifiy password hash match
func (m *Account) Validate() error {
	_, err := govalidator.ValidateStruct(m)
	if err != nil {
		return err
	}
	if m.Password != m.VerifyPassword {
		return errors.New("Model.Account: Password missmatch")
	}
	//er := checkmail.ValidateFormat(m.EmailAccount)
	//
	//if er != nil {
	//	return er
	//}

	m.Email = Email{}
	m.Email.Parse(m.EmailAccount)

	return err

}

// SetPassword create a password hash
func (m *Account) SetPassword(pw string) string {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

// CheckPassword checks that the password hash in the database matches the password the user just gave. Return TRUE if valid
func (m *Account) CheckPassword(dbHash, givenPW string) bool {
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(givenPW)); err != nil {
		return false
	}
	return true
}

func (m *Account) IsValid() error {
	errPers := m.Person.IsValid()
	errComp := m.Company.IsValid(true)
	if (errPers != nil && errComp != nil) || (errPers == nil && errComp == nil) {
		return errors.New("account can't have company and person")
	}

	return nil
}