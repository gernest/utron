# controller

`utron` controllers are structs that implement the `Controller` interface. To help make `utron` usable, `utron` provides a `BaseController` which implements the `Controller` interface and offers additional conveniences to help in composing reusable code.

You get all the benefits of `BaseController` by embedding it in your struct. Our `TODO` Controller is in the `controller/todo.go`

```go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gernest/utron"
	"github.com/gernest/utron/fixtures/todo/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type TODO struct {
	*utron.BaseController
	Routes []string
}

func (t *TODO) Home() {
	todos := []*models.Todo{}
	t.Ctx.DB.Order("created_at desc").Find(&todos)
	t.Ctx.Data["List"] = todos
	t.Ctx.Template = "index"
	t.HTML(http.StatusOK)
}
func (t *TODO) Create() {
	todo := &models.Todo{}
	req := t.Ctx.Request()
	req.ParseForm()
	if err := decoder.Decode(todo, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(todo)
	t.Ctx.Redirect("/", http.StatusFound)
}

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
```

Note that we registered our controller by calling `utron.RegisterController(NewTODO())` in the `init` function
so as to make `utron` aware of our controller. See Routing section below for more explanation of what the controller is doing.

