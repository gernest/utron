package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//PersonName is a controller for PersonName list
type PersonName struct {
	BaseController
	Routes []string
}

//Home renders a PersonName list
func (c *PersonName) Index() {
	PersonNames := []*models.PersonName{}
	c.Ctx.DB.Order("created_at desc").Find(&PersonNames)
	c.Ctx.Data["List"] = PersonNames
	c.Ctx.Template = "application/personname/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a PersonName  item
func (c *PersonName) Create() {
	c.Ctx.Template = "application/personname/index"
	PersonName := &models.PersonName{}

	req := c.Ctx.Request()
	if !c.parseForm(req, PersonName) {
		return
	}

	//Checking that we got valid person name
	if !c.validate(PersonName) {
		return
	}

	rows := c.Ctx.DB.Create(PersonName)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "Can't save person name in database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//View get a PersonName item by id
func (c *PersonName) View() {
	c.Ctx.Template = "application/personname/view"

	PersonNameID := c.Ctx.Params["id"]
	id := c.convertString(PersonNameID)
	if id == -1 {
		return
	}

	PersonName := &models.PersonName{ID: id}
	rows := c.Ctx.DB.Find(PersonName)

	//Checking that this person name is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = PersonName

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

func (c *PersonName) Edit() {
	PersonNameID := c.Ctx.Params["id"]
	id := c.convertString(PersonNameID)
	if id == -1 {
		return
	}

	PersonName := &models.PersonName{ID: id}
	rows := c.Ctx.DB.Find(&PersonName)
	PersonNameFromForm := &models.PersonName{}
	//Checking that this PersonName is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.parseForm(req, PersonNameFromForm) {
		return
	}

	//Checking that we got valid personType
	if !c.validate(PersonNameFromForm) {
		return
	}

	PersonNameFromForm.ID = PersonName.ID
	PersonNameFromForm.CreatedAt = PersonName.CreatedAt
	PersonNameFromForm.UpdatedAt = PersonName.UpdatedAt

	c.Ctx.DB.Save(PersonNameFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//Delete deletes a PersonName item
func (c *PersonName) Delete() {
	PersonNameID := c.Ctx.Params["id"]
	id := c.convertString(PersonNameID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.PersonName{ID: id})

	//Checking that this PersonName was deleted
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/personname", http.StatusFound)
}

//NewPersonName returns a new  PersonName list controller
func NewPersonName() Controller {
	return &PersonName{
		Routes: []string{
			//method;route;handler
			"get;/personname;Index",
			"post;/personname/create;Create",
			"get;/personname/view/{id};View",
			"get;/personname/delete/{id};Delete",
			"post;/personname/update/{id};Edit",
		},
	}
}

func (c *PersonName) validate(personName *models.PersonName) bool {
	err := personName.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *PersonName) isExist(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist person name"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *PersonName) parseForm(req *http.Request, personName *models.PersonName) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(personName, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *PersonName) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
