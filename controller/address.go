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

// Index renders a Address list
func (c *Address) Index() {
	Addresss := []*models.Address{}
	c.Ctx.DB.Order("created_at desc").Find(&Addresss)
	c.Ctx.Data["List"] = Addresss
	c.Ctx.Template = "application/address/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Create creates a Address  item
func (c *Address) Create() {
	c.Ctx.Template = "application/address/index"
	Address := &models.Address{}
	req := c.Ctx.Request()
	if !c.statusInternalServerError(req, Address) {
		return
	}

	//Checking that we got valid address
	if !c.statusBadRequest(Address) {
		return
	}

	rows := c.Ctx.DB.Create(Address)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "can't save Address into database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/address", http.StatusFound)
}

// View a Address item
func (c *Address) View() {
	c.Ctx.Template = "application/address/view"

	AddressID := c.Ctx.Params["id"]
	id := c.convertString(AddressID)
	if id == -1 {
		return
	}

	Address := &models.Address{ID: id}
	rows := c.Ctx.DB.Find(Address)

	//Checking that this address is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Address
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

// Edit allows editing a Address item
func (c *Address) Edit() {
	AddressID := c.Ctx.Params["id"]
	id := c.convertString(AddressID)
	if id == -1 {
		return
	}

	Address := &models.Address{ID: id}
	rows := c.Ctx.DB.Find(&Address)
	AddressFromForm := &models.Address{}
	//Checking that this address is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.statusInternalServerError(req, AddressFromForm) {
		return
	}

	//Checking that we got valid address
	if !c.statusBadRequest(AddressFromForm) {
		return
	}

	AddressFromForm.ID = Address.ID
	AddressFromForm.CreatedAt = Address.CreatedAt
	AddressFromForm.UpdatedAt = Address.UpdatedAt

	c.Ctx.DB.Save(AddressFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/address", http.StatusFound)
}

// Delete deletes a Address item
func (c *Address) Delete() {
	AddressID := c.Ctx.Params["id"]
	id := c.convertString(AddressID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Address{ID: id})

	//Checking that this address was deleted
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/address", http.StatusFound)
}

// NewAddress returns a new  Address controller
func NewAddress() Controller {
	return &Address{
		Routes: []string{
			//method;route;handler
			"get;/address;Index",
			"post;/address/create;Create",
			"get;/address/view/{id};View",
			"get;/address/delete/{id};Delete",
			"post;/address/update/{id};Edit",
		},
	}
}

func (c *Address) statusBadRequest(Address *models.Address) bool {
	err := Address.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Address) statusNotFound(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist address"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Address) statusInternalServerError(req *http.Request, address *models.Address) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(address, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Address) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
