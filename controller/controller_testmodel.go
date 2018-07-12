package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/controller"
	"github.com/NlaakStudios/gowaf/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

//TestModel is a controller for TestModel list
type TestModel struct {
	controller.BaseController
	Routes []string
}

//Home renders a TestModel list
func (t *TestModel) Home() {
	TestModels := []*models.TestModel{}
	t.Ctx.DB.Order("created_at desc").Find(&TestModels)
	t.Ctx.Data["List"] = TestModels
	t.Ctx.Template = "index"
	t.HTML(http.StatusOK)
}

//Create creates a TestModel  item
func (t *TestModel) Create() {
	TestModel := &models.TestModel{}
	req := t.Ctx.Request()
	_ = req.ParseForm()
	if err := decoder.Decode(TestModel, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(TestModel)
	t.Ctx.Redirect("/", http.StatusFound)
}

//Delete deletes a TestModel item
func (t *TestModel) Delete() {
	TestModelID := t.Ctx.Params["id"]
	ID, err := strconv.Atoi(TestModelID)
	if err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}
	t.Ctx.DB.Delete(&models.TestModel{ID: ID})
	t.Ctx.Redirect("/", http.StatusFound)
}

//NewTestModel returns a new  TestModel list controller
func NewTestModel() controller.Controller {
	return &TestModel{
		Routes: []string{
			"get;/;Home",
			"post;/create;Create",
			"get;/delete/{id};Delete",
		},
	}
}
