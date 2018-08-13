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
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Phone  item
func (c *Phone) Create() {
	c.Ctx.Template = "application/phone/index"
	c.Ctx.Data["action"] = "/phone/create"
	Phone := &models.Phone{}

	req := c.Ctx.Request()
	if !c.parseForm(req, Phone) {
		return
	}

	//Checking that we got valid phone
	if !c.validate(Phone) {
		return
	}

	rows := c.Ctx.DB.Create(Phone)

	if rows.RowsAffected != 1 {
		c.Ctx.Data["Message"] = "Can't save phone in database"
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/phone", http.StatusFound)
}

//Get a Phone item by id
func (c *Phone) View() {
	c.Ctx.Template = "application/phone/view"

	PhoneID := c.Ctx.Params["id"]
	id := c.convertString(PhoneID)
	if id == -1 {
		return
	}

	Phone := &models.Phone{ID: id}
	rows := c.Ctx.DB.Find(Phone)

	//Checking that this Phone is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Phone
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

func (c *Phone) Edit() {
	PhoneID := c.Ctx.Params["id"]
	id := c.convertString(PhoneID)
	if id == -1 {
		return
	}

	Phone := &models.Phone{ID: id}
	rows := c.Ctx.DB.Find(&Phone)
	PhoneFromForm := &models.Phone{}

	//Checking that this Phone is exist
	if !c.isExist(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.parseForm(req, PhoneFromForm) {
		return
	}

	//Checking that we got valid Phone
	if !c.validate(PhoneFromForm) {
		return
	}

	PhoneFromForm.ID = Phone.ID
	PhoneFromForm.CreatedAt = Phone.CreatedAt
	PhoneFromForm.UpdatedAt = Phone.UpdatedAt

	c.Ctx.DB.Save(PhoneFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/phone", http.StatusFound)
}

//Delete deletes a Phone item
func (c *Phone) Delete() {
	PhoneID := c.Ctx.Params["id"]
	id := c.convertString(PhoneID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Phone{ID: id})

	//Checking that this Phone was deleted
	if !c.isExist(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/phone", http.StatusFound)
}

//NewPhone returns a new  Phone list controller
func NewPhone() Controller {
	return &Phone{
		Routes: []string{
			//method;route;handler
			"get;/phone;Index",
			"post;/phone/create;Create",
			"get;/phone/view/{id};View",
			"get;/phone/delete/{id};Delete",
			"get;/phone/update/{id};ViewEdit",
			"post;/phone/update/{id};Edit",
		},
	}
}

func (c *Phone) validate(phone *models.Phone) bool {
	err := phone.IsValid()

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusBadRequest)
		return false
	}

	return true
}

func (c *Phone) isExist(rows int64) bool {
	if rows == 0 {
		c.Ctx.Data["Message"] = "Can't manipulate with non exist phone"
		c.Ctx.Template = "error"
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Phone) parseForm(req *http.Request, phone *models.Phone) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(phone, req.PostForm); err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return false
	}

	return true
}

func (c *Phone) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.Data["Message"] = err.Error()
		c.Ctx.Template = "error"
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
