package controller

import (
	"database/sql"
	"errors"
	"fmt"
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
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	emailRoutes = []string{
		"get;/email;Index",
		"post;/email/create;Create",
		"get;/email/view/{id};View",
		"get;/email/delete/{id};Delete",
		"get;/email/update/{id};ViewEdit",
		"post;/email/update/{id};Edit",
	}

	emailAddress = username + "@" + domain
	username     = "someadress"
	domain       = "gmail.com"
	req          *http.Request
	rr           *httptest.ResponseRecorder
	email        *Email
	ctx          *base.Context
	err          error
	mock         sqlmock.Sqlmock

	findQueryEmail   = "SELECT * FROM `emails` WHERE `emails`.`id` = ?"
	deleteQueryEmail = "DELETE FROM `emails` WHERE `emails`.`id` = ?"
	updateQueryEmail = "UPDATE `emails` SET `created_at` = ?, `updated_at` = ?, `address` = ?, `username` = ?, `domain` = ? WHERE `emails`.`id` = ?"

	emailFields = []string{"id", "address", "username", "domain", "created_at", "updated_at"}
	id          = -273
)

func TestEmail_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email")
	email, _ = prepareControllerEmail(req, rr)

	rows := sqlmock.NewRows(emailFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `emails` ORDER BY created_at desc")).WillReturnRows(rows)

	email.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestEmail_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareControllerEmail(req, rr)

	email.prepareValidRequest()
	email.prepareMockRequest()

	email.Create()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create with empty post form
func TestEmail_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareControllerEmail(req, rr)

	req.PostForm = url.Values{}

	email.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create with invalid post form
func TestEmail_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareControllerEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	email.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestEmail_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareControllerEmail(req, rr)

	email.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), emailAddress, username, domain).WillReturnError(errors.New("no connection"))

	email.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist email
func TestEmail_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareControllerEmail(req, rr)

	rows := sqlmock.NewRows(emailFields)
	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnRows(rows.AddRow(id, emailAddress, username, domain, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	email.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get email without param id in request
func TestEmail_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareControllerEmail(req, rr)

	email.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist email
func TestEmail_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareControllerEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	email.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update email with correct data and form
func TestEmail_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareControllerEmail(req, rr)
	newEmailAddress := "sometext" + emailAddress

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)
	req.PostForm = url.Values{}
	req.PostForm.Add("address", newEmailAddress)

	rows := sqlmock.NewRows(emailFields)

	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnRows(rows.AddRow(id, emailAddress, username, domain, time.Now(), time.Now()))

	username = strings.Split(newEmailAddress, "@")[0]
	mock.ExpectExec(fixedFullRe(updateQueryEmail)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), newEmailAddress, username, domain, id).WillReturnResult(sqlmock.NewResult(0, 1))

	email.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to update email without id in params
func TestEmail_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareControllerEmail(req, rr)

	email.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestEmail_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareControllerEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(emailFields)

	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnRows(rows.AddRow(id, emailAddress, username, domain, time.Now(), time.Now()))

	email.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data(email)
func TestEmail_EditInvalidMail(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareControllerEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+domain)

	rows := sqlmock.NewRows(emailFields)

	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnRows(rows.AddRow(id, emailAddress, username, domain, time.Now(), time.Now()))

	email.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist emailAddress
func TestEmail_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareControllerEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	email.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist email
func TestEmail_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareControllerEmail(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryEmail)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	email.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist emailAddress
func TestEmail_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareControllerEmail(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryEmail)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	email.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestEmail_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareControllerEmail(req, rr)

	email.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewEmail(t *testing.T) {
	_ = NewEmail()
}

//Create Request with method and url and ResponseRecorder
func prepareReqAndRecorder(method, url string) (*http.Request, *httptest.ResponseRecorder) {
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	rr = httptest.NewRecorder()

	return req, rr
}

//Create *Email with BaseConroller(with DB) and emailRoutes
func prepareControllerEmail(req *http.Request, rr *httptest.ResponseRecorder) (*Email, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: emailRoutes}
	email := &Email{}

	email.Routes = emailRoutes
	email.BaseController = contr

	return email, ctx
}

func (c *Email) prepareValidRequest() {
	req = c.Ctx.Request()
	req.PostForm = url.Values{}
	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm.Add("address", username+"@"+domain)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *Email) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), emailAddress, username, domain).WillReturnResult(sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}

func fixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}

//func TestEmail_ViewEdit(t *testing.T) {
//	//Create emailAddress with address: someadress@gmail.com
//	req, rr = prepareReqAndRecorder("POST", "/email/create")
//	email, ctx = prepareControllerEmail(req, rr)
//
//	req.PostForm = url.Values{}
//	req.PostForm.Add("address", username+"@"+domain)
//
//	email.Create()
//
//	res := &models.Email{}
//	ctx.DB.Find(res, "address = ?", username+"@"+domain)
//
//	req, rr = prepareReqAndRecorder("GET", "/email/update")
//	email, ctx = prepareControllerEmail(req, rr)
//
//	ctx.Params = make(map[string]string)
//	ctx.Params["id"] = strconv.Itoa(res.ID)
//
//	email.ViewEdit()
//	if rr.Result().StatusCode != http.StatusOK {
//		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
//	}
//
//	ctx.DB.Delete(res)
//
//	//Try to viewEdit without id
//	req, rr = prepareReqAndRecorder("GET", "/email/update")
//	email, ctx = prepareControllerEmail(req, rr)
//
//	email.ViewEdit()
//	if rr.Result().StatusCode != http.StatusInternalServerError {
//		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
//	}
//
//	//Try to viewEdit non exist email
//	req, rr = prepareReqAndRecorder("GET", "/email/update")
//	email, ctx = prepareControllerEmail(req, rr)
//
//	ctx.Params = make(map[string]string)
//	ctx.Params["id"] = strconv.Itoa(123456784)
//
//	email.ViewEdit()
//	if rr.Result().StatusCode != http.StatusNotFound {
//		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
//	}
//}
