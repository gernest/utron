package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Phone is a controller for Phone list
type Phone struct {
	BaseController
	Routes []string
}

//Home renders a Phone list
func (c *Phone) Index() {
	Phones := []*models.Phone{}
	c.Ctx.DB.Order("created_at desc").Find(&Phones)
	c.Ctx.Data["List"] = Phones
	c.Ctx.Template = "application/phone/index"
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.HTML(http.StatusOK)
}

//Create creates a Phone  item
func (c *Phone) Create() {
	Phone := &models.Phone{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Phone, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.DB.Create(Phone)
	c.Ctx.Redirect("/phone", http.StatusFound)
}

//Delete deletes a Phone item
func (c *Phone) Delete() {
	PhoneID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(PhoneID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template, id)
	c.Ctx.DB.Delete(&models.Phone{ID: id})
	c.Ctx.Redirect("/phone", http.StatusFound)
}

//NewPhone returns a new  Phone list controller
func NewPhone() Controller {
	return &Phone{
		Routes: []string{
			"get;/phone;Index",
			"post;/phone/create;Create",
			"get;/phone/delete/{id};Delete",
		},
	}
}
