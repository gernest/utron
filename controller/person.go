package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Person is a controller for Person list
type Person struct {
	BaseController
	Routes []string
}

//Home renders a Person list
func (c *Person) Index() {
	Persons := []*models.Person{}
	c.Ctx.DB.Order("created_at desc").Find(&Persons)
	c.Ctx.Data["List"] = Persons
	c.Ctx.Template = "application/person/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Person  item
func (c *Person) Create() {
	c.Ctx.Template = "application/person/index"
	Person := &models.Person{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Person, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Person)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/person", http.StatusFound)
}

//Delete deletes a Person item
func (c *Person) View() {
	PersonID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.Person{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a Person item
func (c *Person) Delete() {
	PersonID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.Person{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/person", http.StatusFound)
}

//NewPerson returns a new  Person list controller
func NewPerson() Controller {
	return &Person{
		Routes: []string{
			//method;route;handler
			"get;/person;Index",
			"post;/person/create;Create",
			"get;/person/view/{id};View",
			"get;/person/delete/{id};Delete",
		},
	}
}
