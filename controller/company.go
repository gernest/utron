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
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Company  item
func (c *Company) Create() {
	c.Ctx.Template = "application/company/index"
	Company := &models.Company{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(Company, req.PostForm); err != nil {
		c.Ctx.SetError(400, "Internal Server Errror", err.Error())
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(Company)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/company", http.StatusFound)
}

//Delete deletes a Company item
func (c *Company) View() {
	CompanyID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(CompanyID)
	if err != nil {
		c.Ctx.SetError(400, "Internal Server Errror", err.Error())
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Find(&models.Company{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Delete deletes a Company item
func (c *Company) Delete() {
	CompanyID := c.Ctx.Params["id"]
	id, err := strconv.Atoi(CompanyID)
	if err != nil {
		c.Ctx.SetError(400, "Internal Server Errror", err.Error())
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Delete(&models.Company{ID: id})
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/company", http.StatusFound)
}

//NewCompany returns a new  Company list controller
func NewCompany() Controller {
	return &Company{
		Routes: []string{
			//method;route;handler
			"get;/company;Index",
			"post;/company/create;Create",
			"get;/company/view/{id};View",
			"get;/company/delete/{id};Delete",
		},
	}
}
