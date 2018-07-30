package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Gender is a controller for Gender list
type Gender struct {
	BaseController
	Routes []string
}

//Home renders a Gender list
func (c *Gender) Index() {
	Genders := []*models.Gender{}
	c.Ctx.DB.Order("created_at desc").Find(&Genders)
	c.Ctx.Data["List"] = Genders
	c.Ctx.Template = "application/gender/index"
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.HTML(http.StatusOK)
}

//Create creates a Gender  item
func (c *Gender) Create() {
	Gender := &models.Gender{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	c.Ctx.Template = "application/gender/create"
	if err := Decoder.Decode(Gender, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.DB.Create(Gender)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//Delete deletes a Gender item
func (c *Gender) Delete() {
	c.Ctx.Template = "application/gender/delete"
	GenderID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(GenderID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template, id)
	c.Ctx.DB.Delete(&models.Gender{ID: id})
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//NewGender returns a new  Gender list controller
func NewGender() Controller {
	return &Gender{
		Routes: []string{
			"get;/gender;Index",
			"post;/gender/create;Create",
			"get;/gender/delete/{id};Delete",
		},
	}
}
