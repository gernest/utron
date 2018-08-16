package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//PersonType is a controller for PersonType list
type PersonType struct {
	BaseController
	Routes []string
}

//Home renders a PersonType list
func (c *PersonType) Index() {
	PersonTypes := []*models.PersonType{}
	c.Ctx.DB.Order("created_at desc").Find(&PersonTypes)
	c.Ctx.Data["List"] = PersonTypes
	c.Ctx.Template = "application/persontype/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a PersonType  item
func (c *PersonType) Create() {
	c.Ctx.Template = "application/persontype/index"
	PersonType := &models.PersonType{}
	req := c.Ctx.Request()

	if req.Method == "GET" {
		c.Ctx.Template = "application/persontype/create"
		c.Ctx.Data["title"] = "New Person's Type"
		c.Ctx.Data["action"] = "/persontype/create"
		c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
		return
	}

	if !c.parseForm(req, PersonType) {
		return
	}

	//Checking that we got valid personType
	if !c.validate(PersonType) {
		return
	}

	rows := c.Ctx.DB.Create(PersonType)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "Can't save personType in database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/persontype", http.StatusFound)
}

//Get View of PersonType item with id
func (c *PersonType) View() {
	c.Ctx.Template = "application/persontype/view"

	PersonTypeID := c.Ctx.Params["id"]
	id := c.convertString(PersonTypeID)
	if id == -1 {
		return
	}

	PersonType := &models.PersonType{ID: id}
	rows := c.Ctx.DB.Find(PersonType)

	//Checking that this personType is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = PersonType
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

func (c *PersonType) Edit() {
	PersonTypeID := c.Ctx.Params["id"]
	id := c.convertString(PersonTypeID)
	if id == -1 {
		return
	}

	PersonType := &models.PersonType{ID: id}
	rows := c.Ctx.DB.Find(&PersonType)
	PersonTypeFromForm := &models.PersonType{}
	//Checking that this personType is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.parseForm(req, PersonTypeFromForm) {
		return
	}

	//Checking that we got valid personType
	if !c.validate(PersonTypeFromForm) {
		return
	}

	PersonTypeFromForm.ID = PersonType.ID
	PersonTypeFromForm.CreatedAt = PersonType.CreatedAt
	PersonTypeFromForm.UpdatedAt = PersonType.UpdatedAt

	c.Ctx.DB.Save(PersonTypeFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/persontype", http.StatusFound)
}

//Delete deletes a PersonType item
func (c *PersonType) Delete() {
	PersonTypeID := c.Ctx.Params["id"]
	id := c.convertString(PersonTypeID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.PersonType{ID: id})

	//Checking that this personType was deleted
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/persontype", http.StatusFound)
}

//NewPersonType returns a new  PersonType list controller
func NewPersonType() Controller {
	return &PersonType{
		Routes: []string{
			//method;route;handler
			"get;/persontype;Index",
			"get,post;/persontype/create;Create",
			"get;/persontype/view/{id};View",
			"get;/persontype/delete/{id};Delete",
			"post;/persontype/update/{id};Edit",
		},
	}
}

func (c *PersonType) validate(personType *models.PersonType) bool {
	err := personType.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *PersonType) isExist(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist person type"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *PersonType) parseForm(req *http.Request, personType *models.PersonType) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(personType, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *PersonType) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
