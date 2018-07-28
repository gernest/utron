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
func (t *Person) Index() {
	Persons := []*models.Person{}
	t.Ctx.DB.Order("created_at desc").Find(&Persons)
	t.Ctx.Data["List"] = Persons
	t.Ctx.Template = "application/person/index"
	t.HTML(http.StatusOK)
}

//Create creates a Person  item
func (t *Person) Create() {
	Person := &models.Person{}
	req := t.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Person, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(Person)
	t.Ctx.Redirect("/person", http.StatusFound)
}

//Delete deletes a Person item
func (t *Person) Delete() {
	PersonID := t.Ctx.Params["id"]
	id, err := strconv.Atoi(PersonID)
	if err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}
	t.Ctx.DB.Delete(&models.Person{ID: id})
	t.Ctx.Redirect("/person", http.StatusFound)
}

//NewPerson returns a new  Person list controller
func NewPerson() Controller {
	return &Person{
		Routes: []string{
			"get;/person;Index",
			"post;/person/create;Create",
			"get;/person/delete/{id};Delete",
		},
	}
}
