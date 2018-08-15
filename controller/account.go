package controller

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/NlaakStudios/gowaf/models"
	"gopkg.in/gomail.v2"

)

// Account is the controller for the Account Model
type Account struct {
	BaseController
	Routes []string
}

// Index displays the account landing (index) page
func (a *Account) Index() {
	a.Ctx.Template = "application/account/index"
	a.Ctx.Data["title"] = "Home"
	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
}

// Register displays registration page for GET and processes form data on POST
func (a *Account) Register() {
	sess, serr := a.Ctx.SessionStore.New(a.Ctx.Request(), a.Ctx.Cfg.SessionName)
	if serr != nil || sess == nil {
		a.Ctx.Log.Errors("Unable to create new session")
	}

	r := a.Ctx.Request()
	r.ParseForm()
	if r.Method == "GET" {
		a.Ctx.Template = "application/account/register"
		a.Ctx.Data["title"] = "Register"
		a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
		return
	}

	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
	u := &models.Account{}
	err := Decoder.Decode(u, r.PostForm)
	if err != nil {
		//sess.AddFlash("The form was incorrect")
		a.Ctx.Log.Errors(err)
		return
	}

	// Make sure both passwords match
	err = u.Validate()
	if err != nil {
		//sess.AddFlash("The password does not match")
		a.Ctx.Log.Errors(err)
		return
	}

	// Add to database
	u.HashedPassword = u.SetPassword(u.Password)
	u.State = models.UserStateVerifyEmailSent
	a.Ctx.DB.Create(u)

	sess.ID = uuid.New().String()
	sess.Values["uid"] = u.ID
	sess.Values["state"] = u.State
	sess.Values["role"] = u.Access

	//Save session
	errS := a.Ctx.SessionStore.Save(a.Ctx.Request(), a.Ctx.Response(), sess)
	if errS != nil {
		a.Ctx.Log.Errors(errS)
	}

	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
	a.Ctx.Redirect("/account", http.StatusFound)
}

// Login displays login page for GET and processes on POST
func (a *Account) Login() {
	sess, errS := a.Ctx.SessionStore.Get(a.Ctx.Request(), a.Ctx.Cfg.SessionName)
	if errS != nil {
		a.Ctx.Log.Errors("err", errS)
	}

	r := a.Ctx.Request()
	r.ParseForm()
	a.Ctx.Template = "application/account/login"

	if r.Method == "GET" {
		//TODO: Check cookie/session for valid login (ipaddress authroized, etc.) If so use the session to login...
		//else redirect to login page
		if sess.Values["state"] != nil {
			a.Ctx.Data["title"] = "User Login"
			a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
			return
		}
		return
	}

	u := &models.Account{}
	err := Decoder.Decode(u, r.PostForm)
	if err != nil {
		//sess.AddFlash("The form was incorrect")
		a.Ctx.Log.Errors(err)
		return
	}

	var acct models.Account
	db := a.Ctx.DB.First(&acct, "Username = ?", u.Username) // find username with code form username
	//Did we load a a user?
	if acct.ID == 0 {
		a.Ctx.Log.Errors(err)
		//sess.AddFlash("User not found")
		a.Ctx.Redirect("/account/login", http.StatusUnauthorized)
		return
	}

	if acct.CheckPassword(acct.HashedPassword, u.Password) {
		//Login Success - Passwords match
		acct.State = models.UserStateSignedIn
		db.Update(acct)

		sess.ID = uuid.New().String()
		sess.Values["uid"] = u.ID
		sess.Values["state"] = u.State
		sess.Values["role"] = u.Access
		a.Ctx.SessionStore.Save(a.Ctx.Request(), a.Ctx.Response(), sess)
		sess.AddFlash("Login Accepted")
		a.Ctx.Data["loggedin"] = true
		a.Ctx.Template = "application/account/dashboard"
		a.Ctx.Log.Success("Login Accepted")
	} else {
		//Login Success - Passwords match
		a.Ctx.Template = "application/account/login"
		a.Ctx.Log.Errors("Invalid Password")
		a.Ctx.Redirect("/example", http.StatusBadRequest)
		return
	}

	a.Ctx.DB.Model(&models.Account{}).Update(acct)
	a.HTML(http.StatusOK)
}

// Logout logs the user out of they are logged in
func (a *Account) Logout() {
	sess, _ := a.Ctx.SessionStore.Get(a.Ctx.Request(),a.Ctx.Cfg.SessionName)

	r := a.Ctx.Request()
	r.ParseForm()
	uid := sess.Values["uid"]
	if uid == nil {
		a.Ctx.Redirect("/", http.StatusUnauthorized)
	}
	sess.Options.MaxAge = -1
	sess.Save(a.Ctx.Request(), a.Ctx.Response())

	var acct models.Account
	a.Ctx.DB.First(&acct, "ID = ?", uid)
	acct.State = models.UserStateSignedOut
	a.Ctx.DB.Save(acct)
	a.Ctx.Data["loggedin"] = false
	a.Ctx.Template = "application/account/index"
	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
}

// SendEmailVerification Sends a Verification Email to the user registering
func (a *Account) SendEmailVerification(acct models.Account) {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("noreply@%s", a.Ctx.Cfg.Domain))
	m.SetHeader("To", acct.Email)
	m.SetHeader("Subject", fmt.Sprintf("%s email verification", a.Ctx.Cfg.AppName))
	m.SetBody("text/html", fmt.Sprintf("Hello <b>%s</b>, please verify your email address by clicking <a href=\"%s\">here</a>.", acct.Username, a.Ctx.Cfg.BaseURL))

	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

	// Send the email to registering user.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	//Update the account state to Email Verification Sent
	acct.State = models.UserStateVerifyEmailSent
	a.Ctx.DB.Model(&models.Account{}).Update(acct)
	a.HTML(http.StatusOK)
}

// NewAccount returns a new account controller object
func NewAccount() Controller {
	return &Account{
		Routes: []string{
			"get;/account;Index",
			"get,post;/account/register;Register",
			"get,post;/account/login;Login",
			"get;/account/logout;Logout",
		},
	}
}
