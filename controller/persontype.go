package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//PersonType is a controller for PersonType list
type PersonType struct {
	BaseController
	Routes []string
}

//Home renders a PersonType list
func (c *PersonType) Index() {
	PersonTypes := []*models.PersonType{}
	c.Ctx.DB.Order("created_at desc").Find(&PersonTypes)
	c.Ctx.Data["List"] = PersonTypes
	c.Ctx.Template = "application/persontype/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a PersonType  item
func (c *PersonType) Create() {
	c.Ctx.Template = "application/persontype/index"
	PersonType := &models.PersonType{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(PersonType, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(PersonType)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/persontype", http.StatusFound)
}

//Delete deletes a PersonType item
func (c *PersonType) View() {
	PersonTypeID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonTypeID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.PersonType{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a PersonType item
func (c *PersonType) Delete() {
	PersonTypeID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonTypeID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.PersonType{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/persontype", http.StatusFound)
}

//NewPersonType returns a new  PersonType list controller
func NewPersonType() Controller {
	return &PersonType{
		Routes: []string{
			//method;route;handler
			"get;/persontype;Index",
			"post;/persontype/create;Create",
			"get;/persontype/view/{id};View",
			"get;/persontype/delete/{id};Delete",
		},
	}
}
