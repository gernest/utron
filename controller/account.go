package controller

import (
	"net/http"

	"github.com/NlaakStudios/gowaf/models"
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
		// set flash messages
		a.Ctx.Log.Errors(err)
		return
	}

	// Make sure both passwords match
	err = u.Validate()
	if err != nil {
		// set flash messages
		a.Ctx.Log.Errors(err)
		return
	}

	// Add to database
	u.HashedPassword = u.SetPassword(u.Password)
	a.Ctx.DB.Create(u)

	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)

	a.Ctx.Redirect("/account", http.StatusFound)
}

// Login displays login page for GET and processes on POST
func (a *Account) Login() {
	r := a.Ctx.Request()
	r.ParseForm()
	a.Ctx.Template = "application/account/login"
	if r.Method == "GET" {
		a.Ctx.Data["title"] = "User Login"
		a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
		return
	}

	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
	u := &models.Account{}
	err := Decoder.Decode(u, r.PostForm)
	if err != nil {
		// set flash messages
		a.Ctx.Log.Errors(err)
		return
	}

	var acct models.Account
	a.Ctx.DB.First(&acct, "Username = ?", u.Username) // find username with code form username
	if acct.CheckPassword(acct.HashedPassword, u.Password) {
		//Login Success - Passwords match
		acct.State = models.UserStateSignedIn
		a.Ctx.Log.Success("Login Accepted")
		//pretty.Println(acct)
	} else {
		//Login Success - Passwords match
		a.Ctx.Log.Errors("Invalid Password")
	}
}

// Logout logs the user out of they are logged in
func (a *Account) Logout() {
	r := a.Ctx.Request()
	r.ParseForm()
	a.Ctx.Template = "application/account/logout"
	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
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
