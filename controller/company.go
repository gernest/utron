package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NlaakStudios/gowaf/models"
)

// Account is the controller for the Account Model
type Company struct {
	BaseController
	Routes []string
}

// Index renders a Company list
func (c *Company) Index() {
	Companys := []*models.Company{}
	c.Ctx.DB.Order("created_at desc").Find(&Companys)
	c.Ctx.Data["List"] = Companys
	c.Ctx.Template = "application/company/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Create creates a Company  item
func (c *Company) Create() {
	c.Ctx.Template = "application/company/index"
	Company := &models.Company{}
	req := c.Ctx.Request()

	if req.Method == "GET" {
		c.Ctx.Template = "application/company/create"
		c.Ctx.Data["title"] = "New Company"
		c.Ctx.Data["action"] = "/company/create"
		c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
		return
	}

	if !c.statusInternalServerError(req, Company) {
		return
	}

	//Checking that we got valid company
	if !c.statusBadRequest(Company) {
		return
	}

	rows := c.Ctx.DB.First(Company)
	if rows.RowsAffected == 0 {
		//Add New Company
		rows = c.Ctx.DB.Create(Company)
		if rows.RowsAffected != 1 {
			c.Ctx.Data["Message"] = "can't save Company into database"
			c.Ctx.Template = "error"
			c.HTML(http.StatusInternalServerError)
			return
		}
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/company", http.StatusFound)
}

// View a Company item
func (c *Company) View() {
	c.Ctx.Template = "application/company/view"

	CompanyID := c.Ctx.Params["id"]
	id := c.convertString(CompanyID)
	if id == -1 {
		return
	}

	Company := &models.Company{ID: id}
	rows := c.Ctx.DB.Find(Company)

	//Checking that this company is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Company
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Edit allows editing a Company item
func (c *Company) Edit() {
	req := c.Ctx.Request()

	CompanyID := c.Ctx.Params["id"]
	id := c.convertString(CompanyID)
	if id == -1 {
		return
	}

	Company := &models.Company{ID: id}
	rows := c.Ctx.DB.Find(&Company)

	//Checking that this Company is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}
	if req.Method == "GET" {
		c.Ctx.Template = "application/company/create"
		c.Ctx.Data["title"] = "Edit Company"
		c.Ctx.Data["action"] = fmt.Sprintf("/company/update/%d", Company.ID)
		c.Ctx.Data["Payload"] = Company
		c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
		return
	}
}

// Delete deletes a Company item
func (c *Company) Delete() {
	CompanyID := c.Ctx.Params["id"]
	id := c.convertString(CompanyID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Company{ID: id})

	//Checking that this company was deleted
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/company", http.StatusFound)
}

// NewCompany returns a new  Company controller
func NewCompany() Controller {
	return &Company{
		Routes: []string{
			//method;route;handler
			"get;/company;Index",
			"get,post;/company/create;Create",
			"get;/company/view/{id};View",
			"get;/company/delete/{id};Delete",
			"get;/company/edit/{id};Edit",
			"post;/company/update/{id};Update",
		},
	}
}

func (c *Company) statusBadRequest(Company *models.Company) bool {
	err := Company.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Company) statusNotFound(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist company"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Company) statusInternalServerError(req *http.Request, company *models.Company) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(company, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Company) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
