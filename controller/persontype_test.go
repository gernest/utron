package controller

import (
	"github.com/NlaakStudios/gowaf/models"
	"database/sql"
	"log"
	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/logger"
	"os"
	"time"
	"net/url"
	"net/http"
	"testing"
	"strconv"
	"net/http/httptest"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/jinzhu/gorm"
	"errors"
	)

var (
	personTypeRoutes = []string{
		"get;/persontype;Index",
		"post;/persontype/create;Create",
		"get;/persontype/view/{id};View",
		"get;/persontype/delete/{id};Delete",
		"post;/persontype/update/{id};Edit",
	}
	personType *PersonType
	name = "somename"

	findQueryPersonType    = "SELECT * FROM `person_types` WHERE `person_types`.`id` = ?"
	deleteQueryPersonType   = "DELETE FROM `person_types` WHERE `person_types`.`id` = ?"
	updateQueryPersonType  = "UPDATE `person_types` SET `created_at` = ?, `updated_at` = ?, `name` = ? WHERE `person_types`.`id` = ?"

	personTypeFields = []string{"id", "name", "created_at", "updated_at"}
)


func TestPersonType_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes")
	personType, ctx = preparePersonType(req, rr)

	rows := sqlmock.NewRows(personTypeFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `person_types` ORDER BY created_at desc")).WillReturnRows(rows)

	personType.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestPersonType_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/create")
	personType, ctx = preparePersonType(req, rr)

	personType.prepareValidRequest()
	personType.prepareMockRequest()

	personType.Create()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create person type with empty post form
func TestPersonType_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/create")
	personType, ctx = preparePersonType(req, rr)

	req.PostForm = url.Values{}

	personType.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create person type with no valid form params
func TestPersonType_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/create")
	personType, ctx = preparePersonType(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	personType.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestPersonType_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/create")
	personType, ctx = preparePersonType(req, rr)

	personType.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `person_types`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), name).WillReturnError(errors.New("no connection"))

	personType.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist person type
func TestPersonTypeq_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/view")
	personType, ctx = preparePersonType(req, rr)

	rows := sqlmock.NewRows(personTypeFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnRows(rows.AddRow(id, name, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personType.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get person type without param id in request
func TestPersonType_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/view")
	personType, ctx = preparePersonType(req, rr)

	personType.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist person type
func TestPersonType_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/view")
	personType, ctx = preparePersonType(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	personType.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update exist person type
func TestPersonType_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/update")
	personType, ctx = preparePersonType(req, rr)

	newName := "someNewName"

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)
	personType.prepareValidRequest()
	req.PostForm.Add("name", newName)

	rows := sqlmock.NewRows(personTypeFields)

	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnRows(rows.AddRow(id, name, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryPersonType)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), newName, id).WillReturnResult(sqlmock.NewResult(-273, 1))

	personType.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}


//Try to update person type without id in params
func TestPersonType_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/update")
	personType, ctx = preparePersonType(req, rr)

	personType.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestPersonType_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/update")
	personType, ctx = preparePersonType(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(personTypeFields)

	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnRows(rows.AddRow(id, name, time.Now(), time.Now()))

	personType.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestPersonType_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/update")
	personType, ctx = preparePersonType(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("name", "")

	rows := sqlmock.NewRows(personTypeFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnRows(rows.AddRow(id, name, time.Now(), time.Now()))

	personType.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist person type
func TestPersonType_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/persontypes/update")
	personType, ctx = preparePersonType(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	personType.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist person type
func TestPersonType_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/delete")
	personType, ctx = preparePersonType(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPersonType)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personType.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist person type
func TestPersonType_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/delete")
	personType, ctx = preparePersonType(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPersonType)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	personType.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestPersonType_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/persontypes/delete")
	personType, ctx = preparePersonType(req, rr)

	personType.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewPersonType(t *testing.T) {
	_ = NewPersonType()
}

func preparePersonType(req *http.Request, rr *httptest.ResponseRecorder) (*PersonType, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: personTypeRoutes}
	personType := &PersonType{}

	personType.Routes = personTypeRoutes
	personType.BaseController = contr

	return personType, ctx
}

func (c *PersonType) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("name", name)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *PersonType) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `person_types`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), name).WillReturnResult(
		sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}

