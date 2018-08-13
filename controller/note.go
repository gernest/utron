package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Note is a controller for Note list
type Note struct {
	BaseController
	Routes []string
}

//Home renders a Note list
func (c *Note) Index() {
	Notes := []*models.Note{}
	c.Ctx.DB.Order("created_at desc").Find(&Notes)
	c.Ctx.Data["List"] = Notes
	c.Ctx.Template = "application/note/index"
	c.HTML(http.StatusOK)
}

//Create creates a Note  item
func (c *Note) Create() {
	c.Ctx.Template = "application/note/index"
	c.Ctx.Data["action"] = "/note/create"
	Note := &models.Note{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Note, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Note)
	c.Ctx.Redirect("/note", http.StatusFound)
}

//Delete deletes a Note item
func (c *Note) Delete() {
	NoteID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(NoteID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.DB.Delete(&models.Note{ID: id})
	c.Ctx.Redirect("/note", http.StatusFound)
}

//NewNote returns a new  Note list controller
func NewNote() Controller {
	return &Note{
		Routes: []string{
			"get;/note;Index",
			"post;/note/create;Create",
			"get;/note/view/{id};View",
			"get;/note/delete/{id};Delete",
		},
	}
}
