package controller

import (
	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"database/sql"
	"log"
	"github.com/jinzhu/gorm"
	"errors"
)

var (
	genderRoutes = []string{
		"get;/gender;Index",
		"post;/gender/create;Create",
		"get;/gender/view/{id};View",
		"get;/gender/delete/{id};Delete",
	}
	gender     *Gender
	claimedSex = "whatisit"
	bioSex     = byte(0)

	findQueryGender   = "SELECT * FROM `genders` WHERE `genders`.`id` = ?"
	deleteQueryGender = "DELETE FROM `genders` WHERE `genders`.`id` = ?"
	updateQueryGender = "UPDATE `genders` SET `created_at` = ?, `updated_at` = ?, `claimed_sex` = ?, `bio_sex` = ? WHERE `genders`.`id` = ?"

	genderFields = []string{"id", "claimed_sex", "legal_sex", "created_at", "updated_at"}
)

func TestGender_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender")
	gender, _ = prepareGender(req, rr)

	rows := sqlmock.NewRows(genderFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `genders` ORDER BY created_at desc")).WillReturnRows(rows)

	gender.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestGender_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	gender.prepareValidRequest()
	gender.prepareMockRequest()

	gender.Create()

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create gender with empty post form
func TestGender_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	req.PostForm = url.Values{}

	gender.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create gender with no valid form params
func TestGender_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	gender.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestGender_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	gender.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `genders`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), claimedSex, bioSex).WillReturnError(errors.New("no connection"))

	gender.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist gender
func TestGender_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	rows := sqlmock.NewRows(genderFields)
	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnRows(rows.AddRow(id, claimedSex, bioSex, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	gender.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get gender without param id in request
func TestGender_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	gender.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist gender
func TestGender_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	gender.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update gender with correct data and form
func TestGender_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	newClaimedSex := "newClaimedSex"

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)
	req.PostForm = url.Values{}
	req.PostForm.Add("claimed_sex", newClaimedSex)

	rows := sqlmock.NewRows(genderFields)

	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnRows(rows.AddRow(id, claimedSex, bioSex, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryGender)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), newClaimedSex, bioSex, id).WillReturnResult(sqlmock.NewResult(0, 1))

	gender.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}


//Try to update gender without id in params
func TestGender_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	gender.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestGender_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(genderFields)

	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnRows(rows.AddRow(id, claimedSex, bioSex, time.Now(), time.Now()))

	gender.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestGender_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("claimed_sex", "")

	rows := sqlmock.NewRows(genderFields)

	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnRows(rows.AddRow(id, claimedSex, bioSex, time.Now(), time.Now()))

	gender.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist gender
func TestGender_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	gender.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}


//Try to delete exist gender
func TestGender_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/delete")
	gender, ctx = prepareGender(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryGender)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	gender.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist gender
func TestGender_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/delete")
	gender, ctx = prepareGender(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryGender)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	gender.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestGender_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender/delete")
	gender, ctx = prepareGender(req, rr)

	gender.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewGender(t *testing.T) {
	_ = NewGender()
}

func prepareGender(req *http.Request, rr *httptest.ResponseRecorder) (*Gender, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: genderRoutes}
	gender := &Gender{}

	gender.Routes = genderRoutes
	gender.BaseController = contr

	return gender, ctx
}

func (c *Gender) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("claimed_sex", claimedSex)
	req.PostForm.Add("legal_sex", "0")
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *Gender) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `genders`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), claimedSex, bioSex).WillReturnResult(sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}
