package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
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
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.HTML(http.StatusOK)
}

//Create creates a Email  item
func (c *Email) Create() {
	Email := &models.Email{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	c.Ctx.Template = "application/email/create"
	if err := Decoder.Decode(Email, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.DB.Create(Email)
	c.Ctx.Redirect("/email", http.StatusFound)
}

//Delete deletes a Email item
func (c *Email) Delete() {
	c.Ctx.Template = "application/email/delete"
	EmailID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(EmailID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template, id)
	c.Ctx.DB.Delete(&models.Email{ID: id})
	c.Ctx.Redirect("/email", http.StatusFound)
}

//NewEmail returns a new  Email list controller
func NewEmail() Controller {
	return &Email{
		Routes: []string{
			"get;/email;Index",
			"post;/email/create;Create",
			"get;/email/delete/{id};Delete",
		},
	}
}
