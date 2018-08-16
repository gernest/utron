package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

// Account is the controller for the Account Model
type Example struct {
	BaseController
	Routes []string
}

// Index renders a Example list
func (c *Example) Index() {
	Examples := []*models.Example{}
	c.Ctx.DB.Order("created_at desc").Find(&Examples)
	c.Ctx.Data["List"] = Examples
	c.Ctx.Template = "application/example/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Create creates a Example  item
func (c *Example) Create() {
	c.Ctx.Template = "application/example/create"
	c.Ctx.Data["action"] = "/example/create"
	Example := &models.Example{}
	req := c.Ctx.Request()

	if req.Method == "GET" {
		c.Ctx.Data["title"] = "Create Example"
		c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
		return
	}

	if !c.statusInternalServerError(req, Example) {
		return
	}

	//Checking that we got valid Example
	if !c.statusBadRequest(Example) {
		return
	}

	rows := c.Ctx.DB.Create(Example)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "can't save Example into database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/example", http.StatusFound)
}

// View model properies
func (c *Example) View() {
	c.Ctx.Template = "application/example/view"

	ExampleID := c.Ctx.Params["id"]
	id := c.convertString(ExampleID)
	if id == -1 {
		return
	}

	Example := &models.Example{ID: id}
	rows := c.Ctx.DB.Find(Example)

	//Checking that this Example is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Example
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Edit allows editing a Example item
func (c *Example) Edit() {
	req := c.Ctx.Request()

	ExampleID := c.Ctx.Params["id"]
	id := c.convertString(ExampleID)
	if id == -1 {
		return
	}

	Example := &models.Example{ID: id}
	rows := c.Ctx.DB.Find(&Example)

	//Checking that this Example is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}
	if req.Method == "GET" {
		c.Ctx.Template = "application/example/create"
		c.Ctx.Data["title"] = "Edit Example"
		c.Ctx.Data["action"] = fmt.Sprintf("/example/update/%d", Example.ID)
		c.Ctx.Data["Payload"] = Example
		c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
		return
	}
}

// Update allows editing a Example item
func (c *Example) Update() {
	req := c.Ctx.Request()

	ExampleID := c.Ctx.Params["id"]
	id := c.convertString(ExampleID)
	if id == -1 {
		return
	}

	Example := &models.Example{ID: id}
	rows := c.Ctx.DB.Find(&Example)
	ExampleFromForm := &models.Example{}
	//Checking that this Example is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	if !c.statusInternalServerError(req, ExampleFromForm) {
		return
	}

	//Checking that we got valid Example
	if !c.statusBadRequest(ExampleFromForm) {
		return
	}

	ExampleFromForm.ID = Example.ID
	ExampleFromForm.CreatedAt = Example.CreatedAt
	ExampleFromForm.UpdatedAt = Example.UpdatedAt

	c.Ctx.DB.Save(ExampleFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/example", http.StatusFound)
}

// Delete deletes a Example item
func (c *Example) Delete() {
	ExampleID := c.Ctx.Params["id"]
	id := c.convertString(ExampleID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Example{ID: id})

	//Checking that this Example was deleted
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/example", http.StatusFound)
}

// NewExample returns a new  Example controller
func NewExample() Controller {
	return &Example{
		Routes: []string{
			//method;route;handler
			"get;/example;Index",
			"get,post;/example/create;Create",
			"get;/example/view/{id};View",
			"get;/example/delete/{id};Delete",
			"get;/example/edit/{id};Edit",
			"post;/example/update/{id};Update",
		},
	}
}

func (c *Example) statusBadRequest(Example *models.Example) bool {
	err := Example.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Example) statusNotFound(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist Example"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Example) statusInternalServerError(req *http.Request, Example *models.Example) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(Example, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Example) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
