package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

// Example is the controller for the Example Model
type Example struct {
	BaseController
	Routes []string
}

// Index displays the account example (index) page
func (a *Example) Index() {
	a.List()
}

// List shows a paginated list of all model items based on filter / search info
func (a *Example) List() {
	Examples := []*models.Example{}
	a.Ctx.DB.Order("created_at desc").Find(&Examples)
	a.Ctx.Data["List"] = Examples
	a.Ctx.Data["use_styles"] = false
	a.Ctx.Data["use_sparkline"] = false
	a.Ctx.Data["use_datatables"] = true
	a.Ctx.Template = "application/example/index"
	a.HTML(http.StatusOK)
}

// Create creates a new model in the database
func (a *Example) Create() {
	Example := &models.Example{}
	req := a.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Example, req.PostForm); err != nil {
		a.Ctx.Data["Message"] = err.Error()
		a.Ctx.Template = "error"
		a.HTML(http.StatusInternalServerError)
		return
	}

	a.Ctx.DB.Create(Example)
	a.Ctx.Redirect("/example", http.StatusFound)
}


// Delete deletes a model in the database with correct access level
func (a *Example) Delete() {
	ExampleID := a.Ctx.Params["id"]
	id, err := strconv.Atoi(ExampleID)
	if err != nil {
		a.Ctx.Data["Message"] = err.Error()
		a.Ctx.Template = "error"
		a.HTML(http.StatusInternalServerError)
		return
	}
	a.Ctx.DB.Delete(&models.Example{ID: id})
	a.Ctx.Redirect("/example", http.StatusFound)
}

// NewExample returns a new account controller object
func NewExample() Controller {
	return &Example{
		Routes: []string{
			"get;/example;Index",
			"get,post;/example/create;Create",
			"get,post;/example/view/{id};ViewEdit",
			"get;/example/delete/{id};Delete",
		},
	}
}
