package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//PersonName is a controller for PersonName list
type PersonName struct {
	BaseController
	Routes []string
}

//Home renders a PersonName list
func (c *PersonName) Index() {
	PersonNames := []*models.PersonName{}
	c.Ctx.DB.Order("created_at desc").Find(&PersonNames)
	c.Ctx.Data["List"] = PersonNames
	c.Ctx.Template = "application/personname/index"
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.HTML(http.StatusOK)
}

//Create creates a PersonName  item
func (c *PersonName) Create() {
	PersonName := &models.PersonName{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	c.Ctx.Template = "application/personname/create"
	if err := Decoder.Decode(PersonName, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.DB.Create(PersonName)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//Delete deletes a PersonName item
func (c *PersonName) Delete() {
	c.Ctx.Template = "application/personname/delete"
	PersonNameID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonNameID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template, id)
	c.Ctx.DB.Delete(&models.PersonName{ID: id})
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//NewPersonName returns a new  PersonName list controller
func NewPersonName() Controller {
	return &PersonName{
		Routes: []string{
			"get;/personname;Index",
			"post;/personname/create;Create",
			"get;/personname/delete/{id};Delete",
		},
	}
}
