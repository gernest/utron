package controller

import (
	"database/sql"
	"errors"
	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	personNameRoutes = []string{
		"get;/personname;Index",
		"post;/personname/create;Create",
		"get;/personname/view/{id};View",
		"get;/personname/delete/{id};Delete",
		"post;/personname/update/{id};Edit",
	}
	personName *PersonName
	firstName  = "some first name"
	lastName   = "some last name"

	findQueryPersonName   = "SELECT * FROM `person_names` WHERE `person_names`.`id` = ?"
	deleteQueryPersonName = "DELETE FROM `person_names` WHERE `person_names`.`id` = ?"
	updateQueryPersonName = "UPDATE `person_names` SET `created_at` = ?, `updated_at` = ?, `prefix` = ?, `first` = ?, `middle` = ?, `last` = ?, `suffix` = ?, `goes_by` = ? WHERE `person_names`.`id` = ?"

	personNameFields = []string{"id", "first", "last", "created_at", "updated_at"}
)

func TestPersonName_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname")
	personName, ctx = preparePersonName(req, rr)

	rows := sqlmock.NewRows(personNameFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `person_names` ORDER BY created_at desc")).WillReturnRows(rows)

	personName.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestPersonName_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/create")
	personName, ctx = preparePersonName(req, rr)

	personName.prepareValidRequest()
	personName.prepareMockRequest()

	personName.Create()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create person name with empty post form
func TestPersonName_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/create")
	personName, ctx = preparePersonName(req, rr)

	req.PostForm = url.Values{}

	personName.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create person name with no valid form params
func TestPersonName_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/create")
	personName, ctx = preparePersonName(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	personName.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestPersonName_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/create")
	personName, ctx = preparePersonName(req, rr)

	personName.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `person_names`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		firstName, sqlmock.AnyArg(), lastName, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("no connection"))

	personName.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist person name
func TestPersonName_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/view")
	personName, ctx = preparePersonName(req, rr)

	rows := sqlmock.NewRows(personNameFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnRows(rows.AddRow(id, firstName, lastName, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personName.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get person name without param id in request
func TestPersonName_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/view")
	personName, ctx = preparePersonName(req, rr)

	personName.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist person type
func TestPersonName_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/view")
	personName, ctx = preparePersonName(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	personName.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update exist person name
func TestPersonName_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/update")
	personName, ctx = preparePersonName(req, rr)

	newName := "someNewName"

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)
	personName.prepareValidRequest()
	req.PostForm.Add("first", newName)

	rows := sqlmock.NewRows(personNameFields)

	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnRows(rows.AddRow(id, firstName, lastName, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryPersonName)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), newName, sqlmock.AnyArg(), lastName, sqlmock.AnyArg(), sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(-273, 1))

	personName.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to update person name without id in params
func TestPersonName_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/update")
	personName, ctx = preparePersonName(req, rr)

	personName.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestPersonName_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/update")
	personName, ctx = preparePersonName(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(personNameFields)

	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnRows(rows.AddRow(id, firstName, lastName, time.Now(), time.Now()))

	personName.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestPersonName_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/update")
	personName, ctx = preparePersonName(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("first", "")

	rows := sqlmock.NewRows(personNameFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnRows(rows.AddRow(id, firstName, lastName, time.Now(), time.Now()))

	personName.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist person name
func TestPersonName_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/personname/update")
	personName, ctx = preparePersonName(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	personName.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist person name
func TestPersonName_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/delete")
	personName, ctx = preparePersonName(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPersonName)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personName.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist person name
func TestPersonName_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/delete")
	personName, ctx = preparePersonName(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPersonName)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personName.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestPersonName_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/personname/delete")
	personName, ctx = preparePersonName(req, rr)

	personName.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewPersonName(t *testing.T) {
	_ = NewPersonName()
}

func preparePersonName(req *http.Request, rr *httptest.ResponseRecorder) (*PersonName, *base.Context) {
	model := models.NewModel()

	var db *sql.DB

	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("mysql", db)
	if gerr != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(true)

	model.DB = gormDB

	ctx = base.NewContext(rr, req)
	ctx.DB = model
	ctx.Log = logger.NewDefaultLogger(os.Stdout)

	contr := BaseController{Ctx: ctx, Routes: personNameRoutes}
	personName := &PersonName{}

	personName.Routes = personNameRoutes
	personName.BaseController = contr

	return personName, ctx
}

func (c *PersonName) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("first", firstName)
	req.PostForm.Add("last", lastName)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *PersonName) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `person_names`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), firstName, sqlmock.AnyArg(), lastName, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}
