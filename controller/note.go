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
func (t *Note) Index() {
	Notes := []*models.Note{}
	t.Ctx.DB.Order("created_at desc").Find(&Notes)
	t.Ctx.Data["List"] = Notes
	t.Ctx.Template = "application/note/index"
	t.HTML(http.StatusOK)
}

//Create creates a Note  item
func (t *Note) Create() {
	Note := &models.Note{}
	req := t.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Note, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(Note)
	t.Ctx.Redirect("/note", http.StatusFound)
}

//Delete deletes a Note item
func (t *Note) Delete() {
	NoteID := t.Ctx.Params["id"]
	id, err := strconv.Atoi(NoteID)
	if err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}
	t.Ctx.DB.Delete(&models.Note{ID: id})
	t.Ctx.Redirect("/note", http.StatusFound)
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
