package controller

import (
	"net/http"
	"strconv"

	"errors"
	"fmt"

	"github.com/NlaakStudios/gowaf/models"
	"github.com/badoux/checkmail"
)

//Person is a controller for Person list
type Person struct {
	BaseController
	Routes []string
}

//Home renders a Person list
func (c *Person) Index() {
	Persons := []*models.Person{}
	c.Ctx.DB.Preload("Gender").Preload("PersonName").Preload("Email").Preload("PersonType").Preload("Phone").Find(&Persons)
	//c.Ctx.DB.Order("created_at desc").Find(&Persons)
	c.Ctx.Data["List"] = Persons
	c.Ctx.Template = "application/person/index"
	c.HTML(http.StatusOK)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//Create creates a Person  item
func (c *Person) Create() {
	c.Ctx.Template = "application/person/index"

	Person := &models.Person{}

	req := c.Ctx.Request()
	if !c.statusInternalServerError(req, Person) {
		return
	}

	//Checking that we got valid Person
	if !c.statusBadRequest(Person) {
		return
	}

	if err := checkmail.ValidateFormat(Person.Email.Friendly); err == nil {
		emailFromAddress(&Person.Email)
	}

	rows := c.Ctx.DB.Create(Person)

	if rows.RowsAffected != 1 {
		c.Ctx.SetError(http.StatusInternalServerError, "Problem with database", "Can't save Person in database")
		c.HTML(http.StatusInternalServerError)
		return
	}

	if !c.statusBadRequest(Person) {
		c.Ctx.DB.Delete(Person)
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/person", http.StatusFound)
}

//Get View of Person item with id
func (c *Person) View() {
	c.Ctx.Template = "application/person/view"

	PersonID := c.Ctx.Params["id"]
	id := c.convertString(PersonID)
	if id == -1 {
		return
	}

	Person := &models.Person{ID: id}
	rows := c.Ctx.DB.Preload("Gender").Preload("PersonName").Preload("Email").Preload("PersonType").Preload("Phone").Find(Person)

	//Checking that this Person is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Data["Payload"] = Person

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
}

//TODO Now when we update one of nested structure it creates new record for this structure and changes id in person field, i think it's not good, but i don't know how to fix this
func (c *Person) Edit() {
	PersonID := c.Ctx.Params["id"]
	id := c.convertString(PersonID)
	if id == -1 {
		return
	}

	Person := &models.Person{ID: id}
	rows := c.Ctx.DB.Preload("Gender").Preload("PersonName").Preload("Email").Preload("PersonType").Preload("Phone").Find(Person)

	PersonFromForm := &models.Person{}

	//Checking that this person is exist
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	req := c.Ctx.Request()
	if !c.statusInternalServerError(req, PersonFromForm) {
		return
	}

	//Checking that we got valid person
	if !c.statusBadRequest(PersonFromForm) {
		return
	}

	PersonFromForm.ID = Person.ID
	PersonFromForm.CreatedAt = Person.CreatedAt
	PersonFromForm.UpdatedAt = Person.UpdatedAt

	if err := checkmail.ValidateFormat(PersonFromForm.Email.Friendly); err == nil {
		emailFromAddress(&PersonFromForm.Email)
	}

	c.Ctx.DB.Model(&models.Person{}).Update(PersonFromForm)
	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/person", http.StatusFound)
}

//Delete deletes a Person item
func (c *Person) Delete() {
	PersonID := c.Ctx.Params["id"]
	id := c.convertString(PersonID)
	if id == -1 {
		return
	}

	rows := c.Ctx.DB.Delete(&models.Person{ID: id})

	//Checking that this Person was deleted
	if !c.statusNotFound(rows.RowsAffected) {
		return
	}

	c.Ctx.Log.Success(c.Ctx.Request().Method, " : ", c.Ctx.Template)
	c.Ctx.Redirect("/person", http.StatusFound)
}

//NewPerson returns a new  Person list controller
func NewPerson() Controller {
	return &Person{
		Routes: []string{
			//method;route;handler
			"get;/person;Index",
			"post;/person/create;Create",
			"get;/person/view/{id};View",
			"get;/person/delete/{id};Delete",
			"post;/person/update/{id};Edit",
		},
	}
}

func (c *Person) statusBadRequest(person *models.Person) bool {
	var err error

	if person.Dob.IsZero() {
		err = errors.New("invalid date of birthday")
		c.Ctx.SetError(http.StatusBadRequest, "Bad data", err.Error())
		c.HTML(http.StatusBadRequest)
		return false
	}

	if person.PhoneID > 0 {
		phone := &models.Phone{ID: person.PhoneID}
		c.Ctx.DB.Find(phone)
		if err = phone.IsValid(); err != nil {
			c.Ctx.DB.Delete(phone)
			c.Ctx.SetError(http.StatusBadRequest, "Bad data", err.Error())
			c.HTML(http.StatusBadRequest)
			return false
		}
	}

	if person.TypeID > 0 {
		personType := &models.PersonType{ID: person.TypeID}
		c.Ctx.DB.Find(personType)
		if err = personType.IsValid(); err != nil {
			c.Ctx.DB.Delete(personType)
			c.Ctx.SetError(http.StatusBadRequest, "Bad data", err.Error())
			c.HTML(http.StatusBadRequest)
			return false
		}
	}

	if person.GenderID > 0 {
		gender := &models.Gender{ID: person.GenderID}
		c.Ctx.DB.Find(gender)
		if err = gender.IsValid(); err != nil {
			c.Ctx.DB.Delete(gender)
			c.Ctx.SetError(http.StatusBadRequest, "Bad data", err.Error())
			c.HTML(http.StatusBadRequest)
			return false
		}
	}

	if person.NameID > 0 {
		name := &models.PersonName{ID: person.NameID}
		c.Ctx.DB.Find(name)
		if err = name.IsValid(); err != nil {
			c.Ctx.DB.Delete(name)
			c.Ctx.SetError(http.StatusBadRequest, "Bad data", err.Error())
			c.HTML(http.StatusBadRequest)
			return false
		}
	}

	return true
}

func (c *Person) statusNotFound(rows int64) bool {
	if rows == 0 {
		c.Ctx.SetError(http.StatusNotFound, "Not found", "Can't manipulate with non exist person")
		c.HTML(http.StatusNotFound)
		return false
	}

	return true
}

func (c *Person) statusInternalServerError(req *http.Request, person *models.Person) bool {
	_ = req.ParseForm()

	if err := Decoder.Decode(person, req.PostForm); err != nil {
		c.Ctx.SetError(http.StatusInternalServerError, "Internal server error", err.Error())
		c.HTML(http.StatusInternalServerError)
		return false
	}
	fmt.Println("DOB:", person.Dob)
	return true
}

func (c *Person) convertString(id string) int {
	res, err := strconv.Atoi(id)

	if err != nil {
		c.Ctx.SetError(http.StatusInternalServerError, "Internal server error", err.Error())
		c.HTML(http.StatusInternalServerError)
		return -1
	}

	return res
}
