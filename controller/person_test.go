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
	personRoutes = []string{
		"get;/person;Index",
		"post;/person/create;Create",
		"get;/person/view/{id};View",
		"get;/person/delete/{id};Delete",
		"post;/person/update/{id};Edit",
	}
	person         *Person
	dateOfBirthday = time.Now().Format("2006-01-02T15:04:05Z07:00")
	parsed         time.Time

	findQueryPerson   = "SELECT * FROM `people` WHERE `people`.`id` = ?"
	deleteQueryPerson = "DELETE FROM `people` WHERE `people`.`id` = ?"
	updateQueryPerson = "UPDATE `people` SET `created_at` = ?, `dob` = ?, `id` = ?, `updated_at` = ? WHERE `people`.`id` = ?"
	pID = int64(1)
	pEmail    = pUsername + "@" + pDomain
	pUsername = "someaddress"
	pDomain   = "gmail.com"
	personFields = []string{"id", "dob", "created_at", "updated_at"}
)

func TestPerson_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person")
	person, ctx = preparePerson(req, rr)

	rows := sqlmock.NewRows(personTypeFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `people`")).WillReturnRows(rows)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `people` ORDER BY created_at desc")).WillReturnRows(rows)

	person.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestPerson_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	person.prepareValidRequest()

	mock.ExpectBegin()
	person.prepareMockRequest()
	mock.ExpectCommit()

	person.Create()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create person with empty post form
func TestPerson_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	req.PostForm = url.Values{}

	person.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create person with no valid form params
func TestPerson_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	person.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestPerson_CreateWithNoValidName(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `person_names`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		name, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(pID, 1))
	person.prepareValidRequest()
	person.prepareMockRequest()
	mock.ExpectCommit()

	rows := sqlmock.NewRows(personNameFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(1).WillReturnRows(rows.AddRow(id, firstName, "", time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(deleteQueryPersonName)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	req.PostForm.Add("PersonName.first", name)

	person.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestPerson_CreateWithNoValidGender(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	person.prepareValidRequest()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `genders`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), pID).WillReturnResult(sqlmock.NewResult(pID, 1))
	person.prepareMockRequest()
	mock.ExpectCommit()

	rows := sqlmock.NewRows(genderFields)
	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(pID).WillReturnRows(rows.AddRow(pID, "", 1, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(deleteQueryGender)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))

	req.PostForm.Add("Gender.legal_sex", "1")

	person.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestPerson_CreateWithNoValidPhone(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	person.prepareValidRequest()

	mock.ExpectBegin()
	formatTime()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pEmail,
		pUsername, pDomain).WillReturnResult(sqlmock.NewResult(pID, 1))
	mock.ExpectExec("INSERT INTO `phones`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "1", sqlmock.AnyArg(),
		sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(pID, 1))
	mock.ExpectExec("INSERT INTO `people`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), parsed,
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(pID, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows(phoneFields)
	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(pID).WillReturnRows(rows.AddRow(pID, 1, sqlmock.AnyArg(),
		sqlmock.AnyArg(), 0, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(deleteQueryPhone)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))

	req.PostForm.Add("Phone.code", "1")

	person.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestPerson_CreateWithNoValidType(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	person.prepareValidRequest()

	mock.ExpectBegin()
	formatTime()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pEmail,
		pUsername, pDomain).WillReturnResult(sqlmock.NewResult(pID, 1))
	mock.ExpectExec("INSERT INTO `person_types`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "name").WillReturnResult(
		sqlmock.NewResult(pID, 1))
	mock.ExpectExec("INSERT INTO `people`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), parsed,
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(pID, 1))
	mock.ExpectCommit()

	rows := sqlmock.NewRows(personTypeFields)
	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(pID).WillReturnRows(rows.AddRow(pID, "", time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(deleteQueryPersonType)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(pID).WillReturnResult(sqlmock.NewResult(0, 1))

	req.PostForm.Add("PersonType.name", "name")

	person.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestPerson_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/create")
	person, ctx = preparePerson(req, rr)

	person.prepareValidRequest()
	formatTime()
	mock.ExpectExec("INSERT INTO `people`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), parsed,
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errors.New("no connection"))

	person.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist person
func TestPerson_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/view")
	person, ctx = preparePerson(req, rr)

	rows := sqlmock.NewRows(personFields)
	formatTime()
	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnRows(rows.AddRow(id, parsed, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	person.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get person without param id in request
func TestPerson_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/view")
	person, ctx = preparePerson(req, rr)

	person.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist person
func TestPerson_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/view")
	person, ctx = preparePerson(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	person.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update exist person
func TestPerson_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/update")
	person, ctx = preparePerson(req, rr)

	formatTime()
	newDob := parsed.Add(time.Hour * 1024)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	person.prepareValidRequest()

	req.PostForm.Add("dob", newDob.Format("2006-01-02T15:04:05Z07:00"))

	rows := sqlmock.NewRows(personFields)

	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnRows(rows.AddRow(id, parsed, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryPerson)).WithArgs(sqlmock.AnyArg(), newDob, id, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(-273, 1))

	person.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to update person without id in params
func TestPerson_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/update")
	person, ctx = preparePerson(req, rr)

	person.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestPerson_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/update")
	person, ctx = preparePerson(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(personFields)

	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnRows(rows.AddRow(id, time.Now(), time.Now(), time.Now()))

	person.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestPerson_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/update")
	person, ctx = preparePerson(req, rr)
	invalidID := 10000000000

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("gender_id", strconv.Itoa(invalidID))
	req.PostForm.Add("name_id", strconv.Itoa(invalidID))
	req.PostForm.Add("phone_id", strconv.Itoa(invalidID))
	req.PostForm.Add("email_id", strconv.Itoa(invalidID))
	req.PostForm.Add("type_id", strconv.Itoa(invalidID))

	rows := sqlmock.NewRows(personFields)
	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnRows(rows.AddRow(id, time.Now(), time.Now(), time.Now()))
	mock.ExpectQuery(fixedFullRe(findQueryPhone)).WithArgs(invalidID).WillReturnError(errors.New("doesn't exist"))
	mock.ExpectQuery(fixedFullRe(findQueryPersonType)).WithArgs(invalidID).WillReturnError(errors.New("doesn't exist"))
	mock.ExpectQuery(fixedFullRe(findQueryEmail)).WithArgs(invalidID).WillReturnError(errors.New("doesn't exist"))
	mock.ExpectQuery(fixedFullRe(findQueryGender)).WithArgs(invalidID).WillReturnError(errors.New("doesn't exist"))
	mock.ExpectQuery(fixedFullRe(findQueryPersonName)).WithArgs(invalidID).WillReturnError(errors.New("doesn't exist"))

	person.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist person
func TestPerson_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/person/update")
	person, ctx = preparePerson(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryPerson)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	person.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist person
func TestPerson_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/delete")
	person, ctx = preparePerson(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	person.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist person
func TestPerson_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/delete")
	person, ctx = preparePerson(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryPerson)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	person.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestPerson_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/person/delete")
	person, ctx = preparePerson(req, rr)

	person.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewPerson(t *testing.T) {
	_ = NewPerson()
}

func preparePerson(req *http.Request, rr *httptest.ResponseRecorder) (*Person, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: personRoutes}
	person := &Person{}

	person.Routes = personRoutes
	person.BaseController = contr

	return person, ctx
}

func (c *Person) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("dob", dateOfBirthday)
	req.PostForm.Add("Email.address", pEmail)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *Person) prepareMockRequest() {
	formatTime()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pEmail,
		pUsername, pDomain).WillReturnResult(sqlmock.NewResult(pID, 1))
	mock.ExpectExec("INSERT INTO `people`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), parsed,
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(
		sqlmock.NewResult(pID, 1))
}

func formatTime() {
	parsed, err = time.Parse("2006-01-02T15:04:05Z07:00", dateOfBirthday)

	if err != nil {
		panic(err)
	}
}
