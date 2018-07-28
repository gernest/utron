package controller

import (
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

//Company is a controller for Company list
type Company struct {
	BaseController
	Routes []string
}

//Home renders a Company list
func (c *Company) Index() {
	Companys := []*models.Company{}
	c.Ctx.DB.Order("created_at desc").Find(&Companys)
	c.Ctx.Data["List"] = Companys
	c.Ctx.Template = "application/company/index"
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.HTML(http.StatusOK)
}

//Create creates a Company  item
func (c *Company) Create() {
	Company := &models.Company{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	c.Ctx.Template = "application/company/create"
	if err := Decoder.Decode(Company, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.DB.Create(Company)
	c.Ctx.Redirect("/company", http.StatusFound)
}

//Delete deletes a Company item
func (c *Company) Delete() {
	c.Ctx.Template = "application/company/delete"
	CompanyID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(CompanyID)
	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.Ctx.Log.Errors(err)
		c.HTML(http.StatusInternalServerError)
		return
	}
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template, id)
	c.Ctx.DB.Delete(&models.Company{ID: id})
	c.Ctx.Redirect("/company", http.StatusFound)
}

//NewCompany returns a new  Company list controller
func NewCompany() Controller {
	return &Company{
		Routes: []string{
			"get;/company;Index",
			"post;/company/create;Create",
			"get;/company/delete/{id};Delete",
		},
	}
}
