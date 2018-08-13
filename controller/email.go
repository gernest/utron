package controller

import (
	"net/http"
	"strconv"

	"strings"

	"github.com/NlaakStudios/gowaf/models"
	"github.com/badoux/checkmail"
)

//Email is a controller for Email list
type Email struct {
	BaseController
	Routes []string
}

//Home renders a Email list
func (c *Email) Index() {
	Emails := []*models.Email{}
	c.Ctx.DB.Order("created_at desc").Find(&Emails)
	c.Ctx.Data["List"] = Emails
	c.Ctx.Template = "application/email/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Email  item
func (c *Email) Create() {
	c.Ctx.Template = "application/email/index"

	Email := &models.Email{}

	req := c.Ctx.Request()
	if !c.parseForm(req, Email) {
		return
	}

	//Checking that we got valid emailAddress
	if !c.validate(Email) {
		return
	}

	//Add username and host to email
	emailFromAddress(Email)

	rows := c.Ctx.DB.Create(Email)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "Can't save email in database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/email", http.StatusFound)
}

//Delete deletes a Email item
func (c *Email) View() {
	c.Ctx.Template = "application/email/view"

	EmailID := c.Ctx.Params["id"]
	id := c.convertString(EmailID)

	if id == -1 {
		return
	}

	Email := &models.Email{ID: id}
	rows := c.Ctx.DB.Find(Email)

	//Checking that this email is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Email
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

func (c *Email) Edit() {
	EmailID := c.Ctx.Params["id"]
	id := c.convertString(EmailID)

	if id == -1 {
		return
	}

	Email := &models.Email{ID: id}
	rows := c.Ctx.DB.Find(Email)
	EmailFromForm := &models.Email{}

	//Checking that this email is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.parseForm(req, EmailFromForm) {
		return
	}

	//Checking that we got valid emailAddress
	if !c.validate(EmailFromForm) {
		return
	}

	//Add username and host to email
	Email.Friendly = EmailFromForm.Friendly
	emailFromAddress(Email)

	c.Ctx.DB.Save(Email)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/email", http.StatusFound)
}

//Delete deletes a Email item
func (c *Email) Delete() {
	EmailID := c.Ctx.Params["id"]
	id := c.convertString(EmailID)

	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Email{ID: id})

	//Checking that this email is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/email", http.StatusFound)
}

//NewEmail returns a new  Email list controller
func NewEmail() Controller {
	return &Email{
		Routes: []string{
			//method;route;handler
			"get;/email;Index",
			"post;/email/create;Create",
			"get;/email/view/{id};View",
			"get;/email/delete/{id};Delete",
			//"get;/email/update/{id};ViewEdit",
			"post;/email/update/{id};Edit",
		},
	}
}

func emailFromAddress(email *models.Email) {
	str := strings.Split(email.Friendly, "@")
	email.Username = str[0]
	email.Domain = str[1]
}

func (c *Email) validate(Email *models.Email) bool {
	err := checkmail.ValidateFormat(Email.Friendly)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Email) isExist(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist email"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Email) parseForm(req *http.Request, email *models.Email) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(email, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Email) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
