package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Gender is a controller for Gender list
type Gender struct {
	BaseController
	Routes []string
}

//Home renders a Gender list
func (c *Gender) Index() {
	Genders := []*models.Gender{}
	c.Ctx.DB.Order("created_at desc").Find(&Genders)
	c.Ctx.Data["List"] = Genders
	c.Ctx.Template = "application/gender/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Gender  item
func (c *Gender) Create() {
	c.Ctx.Template = "application/gender/index"
	c.Ctx.Data["action"] = "/gender/create"

	Gender := &models.Gender{}

	req := c.Ctx.Request()
	if !c.parseForm(req, Gender) {
		return
	}

	if !c.validate(Gender) {
		return
	}

	rows := c.Ctx.DB.Create(Gender)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "Can't save gender in database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//Delete deletes a Gender item
func (c *Gender) View() {
	c.Ctx.Template = "application/gender/view"
	GenderID := c.Ctx.Params["id"]
	id := c.convertString(GenderID)
	if id == -1 {
		return
	}

	Gender := &models.Gender{ID: id}
	rows := c.Ctx.DB.Find(Gender)

	//Checking that this gender is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Gender

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

func (c *Gender) Edit() {
	GenderID := c.Ctx.Params["id"]
	id := c.convertString(GenderID)
	if id == -1 {
		return
	}

	Gender := &models.Gender{ID: id}
	rows := c.Ctx.DB.Find(&Gender)
	GenderFromForm := &models.Gender{}

	//Checking that this gender is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.parseForm(req, GenderFromForm) {
		return
	}

	//Checking that we got valid gender
	if !c.validate(GenderFromForm) {
		return
	}

	GenderFromForm.ID = Gender.ID
	GenderFromForm.CreatedAt = Gender.CreatedAt
	GenderFromForm.UpdatedAt = Gender.UpdatedAt

	c.Ctx.DB.Save(GenderFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//TODO
//func (c *Gender) ViewEdit() {
//	c.Ctx.Template = "application/gender/update"
//	GenderID := c.Ctx.Params["id"]
//	id, err := strconv.Atoi(GenderID)
//	if err != nil {
//		c.Ctx.Data["Message"] = err.Error()
//		c.Ctx.Template = "error"
//		c.HTML(http.StatusInternalServerError)
//		return"errors"
//	}
//
//	Gender := &models.Gender{ID: id}
//	c.Ctx.DB.Find(&Gender)
//
//	req := c.Ctx.Request()
//	_ = req.ParseForm()
//	if err := Decoder.Decode(Gender, req.PostForm); err != nil {
//		c.Ctx.Data["Message"] = err.Error()
//		c.Ctx.Template = "error"
//		c.HTML(http.StatusInternalServerError)
//		return
//	}
//
//	c.Ctx.DB.Save(Gender)
//	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
//	c.Ctx.Redirect("/gender", http.StatusFound)
//}

//Delete deletes a Gender item
func (c *Gender) Delete() {
	GenderID := c.Ctx.Params["id"]
	id := c.convertString(GenderID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Gender{ID: id})

	//Checking that this gender was deleted
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/gender", http.StatusFound)
}

//NewGender returns a new  Gender list controller
func NewGender() Controller {
	return &Gender{
		Routes: []string{
			//method;route;handler
			"get;/gender;Index",
			"post;/gender/create;Create",
			"get;/gender/view/{id};View",
			"get;/gender/delete/{id};Delete",
			"post;/gender/update/{id};Edit",
		},
	}
}

func (c *Gender) validate(gender *models.Gender) bool {
	err := gender.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Gender) isExist(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist gender"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Gender) parseForm(req *http.Request, gender *models.Gender) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(gender, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Gender) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
