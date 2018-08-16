package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	phoneRoutes = []string{
		"get;/phone;Index",
		"post;/phone/create;Create",
		"get;/phone/view/{id};View",
		"get;/phone/delete/{id};Delete",
		"get;/phone/update/{id};ViewEdit",
		"post;/phone/update/{id};Edit",
	}
	phone       *Phone
	phoneNumber = "234567"
	countryCode = "380"

	findQueryPhone   = "SELECT * FROM `phones` WHERE `phones`.`id` = ?"
	deleteQueryPhone = "DELETE FROM `phones` WHERE `phones`.`id` = ?"
	updateQueryPhone = "UPDATE `phones` SET `created_at` = ?, `updated_at` = ?, `country_code` = ?, `area_code` = ?, `number` = ?, `phone_type` = ? WHERE `phones`.`id` = ?"

	phoneFields = []string{"id", "code", "area", "number", "phone_type", "created_at", "updated_at"}
)

func TestPhone_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone")
	phone, ctx = preparePhone(req, rr)

	rows := sqlmock.NewRows(phoneFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `phones` ORDER BY created_at desc")).WillReturnRows(rows)

	phone.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestPhone_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/create")
	phone, ctx = preparePhone(req, rr)

	phone.prepareValidRequest()
	phone.prepareMockRequest()

	phone.Create()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create phone with empty post form
func TestPhone_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/create")
	phone, ctx = preparePhone(req, rr)

	req.PostForm = url.Values{}

	phone.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create phone with no valid form params
func TestPhone_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/create")
	phone, ctx = preparePhone(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	phone.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestPhone_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/create")
	phone, ctx = preparePhone(req, rr)

	phone.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `phones`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), countryCode, sqlmock.AnyArg(),
		phoneNumber, sqlmock.AnyArg()).WillReturnError(errors.New("no connection"))

	phone.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist phone
func TestPhone_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/view")
	phone, ctx = preparePhone(req, rr)

	rows := sqlmock.NewRows(phoneFields)
	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnRows(rows.AddRow(id, countryCode, sqlmock.AnyArg(),
		phoneNumber, 0, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	phone.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get phone without param id in request
func TestPhone_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/view")
	phone, ctx = preparePhone(req, rr)

	phone.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist phone
func TestPhone_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/view")
	phone, ctx = preparePhone(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	phone.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update exist phone
func TestPhone_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/update")
	phone, ctx = preparePhone(req, rr)

	newNumber := "+380994"

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	phone.prepareValidRequest()

	req.PostForm.Set("number", newNumber)

	rows := sqlmock.NewRows(phoneFields)

	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnRows(rows.AddRow(id, countryCode, sqlmock.AnyArg(),
		phoneNumber, 0, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryPhone)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), countryCode, sqlmock.AnyArg(),
		newNumber, sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(-273, 1))

	phone.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to update phone without id in params
func TestPhone_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/update")
	phone, ctx = preparePhone(req, rr)

	phone.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestPhone_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/update")
	phone, ctx = preparePhone(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(phoneFields)

	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnRows(rows.AddRow(id, countryCode, sqlmock.AnyArg(),
		phoneNumber, 0, time.Now(), time.Now()))

	phone.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestPhone_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/update")
	phone, ctx = preparePhone(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("number", "")

	rows := sqlmock.NewRows(phoneFields)

	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnRows(rows.AddRow(id, countryCode, sqlmock.AnyArg(),
		phoneNumber, 0, time.Now(), time.Now()))

	phone.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist phone
func TestPhone_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/phone/update")
	phone, ctx = preparePhone(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	phone.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist phone
func TestPhone_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/delete")
	phone, ctx = preparePhone(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPhone)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	phone.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist phone
func TestPhone_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/delete")
	phone, ctx = preparePhone(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPhone)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	phone.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestPhone_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/phone/delete")
	phone, ctx = preparePhone(req, rr)

	phone.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewPhone(t *testing.T) {
	_ = NewPhone()
}

func preparePhone(req *http.Request, rr *httptest.ResponseRecorder) (*Phone, *base.Context) {
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

	ctx := base.NewContext(rr, req)
	ctx.DB = model
	ctx.Log = logger.NewDefaultLogger(os.Stdout)

	contr := BaseController{Ctx: ctx, Routes: phoneRoutes}
	phone := &Phone{}

	phone.Routes = phoneRoutes
	phone.BaseController = contr

	return phone, ctx
}

func (c *Phone) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("code", countryCode)
	req.PostForm.Add("number", phoneNumber)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *Phone) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `phones`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), countryCode, sqlmock.AnyArg(),
		phoneNumber, sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}
