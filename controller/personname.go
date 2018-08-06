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
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a PersonName  item
func (c *PersonName) Create() {
	c.Ctx.Template = "application/personname/index"
	PersonName := &models.PersonName{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(PersonName, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(PersonName)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//Delete deletes a PersonName item
func (c *PersonName) View() {
	PersonNameID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonNameID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.PersonName{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a PersonName item
func (c *PersonName) Delete() {
	PersonNameID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonNameID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.PersonName{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//NewPersonName returns a new  PersonName list controller
func NewPersonName() Controller {
	return &PersonName{
		Routes: []string{
			//method;route;handler
			"get;/personname;Index",
			"post;/personname/create;Create",
			"get;/personname/view/{id};View",
			"get;/personname/delete/{id};Delete",
		},
	}
}
