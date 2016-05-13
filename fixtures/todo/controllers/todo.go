package controllers

import (
	"net/http"
	"strconv"

	"github.com/gernest/utron"
	"github.com/gernest/utron/fixtures/todo/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

//TODO is a controller for Todo list
type TODO struct {
	*utron.BaseController
	Routes []string
}

//Home renders a todo list
func (t *TODO) Home() {
	todos := []*models.Todo{}
	t.Ctx.DB.Order("created_at desc").Find(&todos)
	t.Ctx.Data["List"] = todos
	t.Ctx.Template = "index"
	t.HTML(http.StatusOK)
}

//Create creates a todo  item
func (t *TODO) Create() {
	todo := &models.Todo{}
	req := t.Ctx.Request()
	_ = req.ParseForm()
	if err := decoder.Decode(todo, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(todo)
	t.Ctx.Redirect("/", http.StatusFound)
}

//Delete deletes a todo item
func (t *TODO) Delete() {
	todoID := t.Ctx.Params["id"]
	ID, err := strconv.Atoi(todoID)
	if err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}
	t.Ctx.DB.Delete(&models.Todo{ID: ID})
	t.Ctx.Redirect("/", http.StatusFound)
}

//NewTODO returns a new  todo list controller
func NewTODO() *TODO {
	return &TODO{
		Routes: []string{
			"get;/;Home",
			"post;/create;Create",
			"get;/delete/{id};Delete",
		},
	}
}

func init() {
	utron.RegisterController(NewTODO())
}
