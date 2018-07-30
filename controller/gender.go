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
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Gender  item
func (c *Gender) Create() {
	c.Ctx.Template = "application/gender/index"
	Gender := &models.Gender{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Gender, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Gender)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//Delete deletes a Gender item
func (c *Gender) View() {
	GenderID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(GenderID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.Gender{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a Gender item
func (c *Gender) Delete() {
	GenderID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(GenderID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.Gender{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//NewGender returns a new  Gender list controller
func NewGender() Controller {
	return &Gender{
		Routes: []string{
			//method;route;handler
			"get;/gender;Index",
			"post;/gender/create;Create",
			"get;/gender/view/{id};View",
			"get;/gender/delete/{id};Delete",
		},
	}
}
