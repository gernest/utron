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
	addressRoutes = []string{
		"get;/address;Index",
		"post;/address/create;Create",
		"get;/address/view/{id};View",
		"get;/address/delete/{id};Delete",
		"get;/address/update/{id};ViewEdit",
		"post;/address/update/{id};Edit",
	}
	address *Address
	number  = 10000000
	street  = "any street"
	city    = "any city"
	state   = "any state"
	zip     = "any zip"
	county  = "any county"
	country = "any country"

	findQueryAddresses     = "SELECT * FROM `addresses` WHERE `addresses`.`id` = ?"
	deleteQueryAddresses   = "DELETE FROM `addresses` WHERE `addresses`.`id` = ?"
	updateQueryAddresses   = "UPDATE `addresses` SET `created_at` = ?, `updated_at` = ?, `number` = ?, `street` = ?, `city` = ?, `state` = ?," +
		" `zip` = ?, `county` = ?, `country` = ? WHERE `addresses`.`id` = ?"

	//UPDATE `addresses` SET `created_at` = ?, `updated_at` = ?, `number` = ?, `street` = ?, `city` = ?, `state` = ?, `zip` = ?, `county` = ?, `country` = ? WHERE `addresses`.`id` =

	addressesFields = []string{"id", "number", "street", "city", "state", "zip", "county", "country", "created_at", "updated_at"}
)

func TestAddress_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address")
	address, _ = prepareAddress(req, rr)

	rows := sqlmock.NewRows(addressesFields)
	mock.ExpectQuery(fixedFullRe("SELECT * FROM `addresses` ORDER BY created_at desc")).WillReturnRows(rows)

	address.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to create with valid data
func TestAddress_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)

	address.prepareValidRequest()
	address.prepareMockRequest()

	address.Create()

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to create addresses with empty post form
func TestAddress_CreateWithEmptyForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)

	req.PostForm = url.Values{}

	address.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to create addresses with no valid form params
func TestAddress_CreateWithNoValidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	address.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to create without connection to DB
func TestAddress_CreateWithNoConnToDB(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)

	address.prepareValidRequest()
	mock.ExpectExec("INSERT INTO `addresses`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), number, street,
		city, state, zip, county, country).WillReturnError(errors.New("no connection"))

	address.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get exist addresses
func TestAddressq_View(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	rows := sqlmock.NewRows(addressesFields)
	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnRows(rows.AddRow(id, number, street,
		city, state, zip, county, country, time.Now(), time.Now()))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	address.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

//Try to get addresses without param id in request
func TestAddress_ViewWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	address.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to get no exist addresses
func TestAddress_ViewNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	address.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to update exist address
func TestAddress_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	newNumber := 1

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)
	address.prepareValidRequest()
	req.PostForm.Add("number", strconv.Itoa(newNumber))

	rows := sqlmock.NewRows(addressesFields)

	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnRows(rows.AddRow(id, number, street,
		city, state, zip, county, country, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryAddresses)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), newNumber, street,
		city, state, zip, county, country, id).WillReturnResult(sqlmock.NewResult(-273, 1))

	address.Edit()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}


//Try to update address without id in params
func TestAddress_EditWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	address.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid field
func TestAddress_EditInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	rows := sqlmock.NewRows(addressesFields)

	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnRows(rows.AddRow(id, number, street,
		city, state, zip, county, country, time.Now(), time.Now()))

	address.Edit()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

//Try to post form with invalid data
func TestAddress_EditInvalidData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	req.PostForm = url.Values{}
	req.PostForm.Add("number", strconv.Itoa(111111))
	req.PostForm.Add("street", "dafdsf")

	rows := sqlmock.NewRows(addressesFields)

	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnRows(rows.AddRow(id, number, street,
		city, state, zip, county, country, time.Now(), time.Now()))

	address.Edit()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

//Try to update no exist address
func TestAddress_EditNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	mock.ExpectQuery(fixedFullRe(findQueryAddresses)).WithArgs(id).WillReturnError(errors.New("doesn't exist"))

	address.Edit()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete exist address
func TestAddress_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/delete")
	address, ctx = prepareAddress(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryAddresses)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	address.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//Try to delete no exist address
func TestAddress_DeleteNoExist(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/delete")
	address, ctx = prepareAddress(req, rr)

	mock.ExpectExec(fixedFullRe(deleteQueryAddresses)).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(id)

	address.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//Try to delete without param id
func TestAddress_DeleteWithoutID(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address/delete")
	address, ctx = prepareAddress(req, rr)

	address.Delete()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestNewAddress(t *testing.T) {
	_ = NewAddress()
}

func prepareAddress(req *http.Request, rr *httptest.ResponseRecorder) (*Address, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: addressRoutes}
	address := &Address{}

	address.Routes = addressRoutes
	address.BaseController = contr

	return address, ctx
}

func (c *Address) prepareValidRequest() {
	req = c.Ctx.Request()

	created := time.Now().Format("2006-01-02T15:04:05Z07:00")

	req.PostForm = url.Values{}
	req.PostForm.Add("number", strconv.Itoa(number))
	req.PostForm.Add("street", street)
	req.PostForm.Add("state", state)
	req.PostForm.Add("city", city)
	req.PostForm.Add("zip", zip)
	req.PostForm.Add("county", county)
	req.PostForm.Add("country", country)
	req.PostForm.Add("created", created)
	req.PostForm.Add("updated", created)
}

func (c *Address) prepareMockRequest() {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `addresses`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), number, street,
		city, state, zip, county, country).WillReturnResult(sqlmock.NewResult(-273, 1))
	mock.ExpectCommit()
}

