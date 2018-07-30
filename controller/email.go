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
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Email  item
func (c *Email) Create() {
	c.Ctx.Template = "application/email/index"
	Email := &models.Email{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Email, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Email)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/email", http.StatusFound)
}

//Delete deletes a Email item
func (c *Email) View() {
	EmailID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(EmailID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.Email{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a Email item
func (c *Email) Delete() {
	EmailID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(EmailID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.Email{ID: id})
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
		},
	}
}
