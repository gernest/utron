package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

// Account is the controller for the Account Model
type Address struct {
	BaseController
	Routes []string
}

//Home renders a Address list
func (c *Address) Index() {
	Addresss := []*models.Address{}
	c.Ctx.DB.Order("created_at desc").Find(&Addresss)
	c.Ctx.Data["List"] = Addresss
	c.Ctx.Template = "application/address/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Address  item
func (c *Address) Create() {
	Address := &models.Address{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Address, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Address)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/address", http.StatusFound)
}

//Delete deletes a Address item
func (c *Address) Delete() {
	AddressID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(AddressID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}
	//TODO: How to compare gorm.Model.ID
	c.Ctx.DB.Delete(&models.Address{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/address", http.StatusFound)
}

//NewAddress returns a new  Address list controller
func NewAddress() Controller {
	return &Address{
		Routes: []string{
			"get;/address;Index",
			"post;/address/create;Create",
			"get;/address/delete/{id};Delete",
		},
	}
}
