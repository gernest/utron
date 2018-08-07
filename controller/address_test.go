package controller

import (
	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/config"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
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
)

func TestAddress_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/address")
	address, _ = prepareAddress(req, rr)

	address.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestAddress_Create(t *testing.T) {
	//Try to create address with invalid form params
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	address, ctx = prepareAddress(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	address.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Create address with params
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)
	prepareRequest(req)

	address.Create()

	res := &models.Address{}
	ctx.DB.Find(res, "number = ?", number)

	if res.Number != number || res.Street != street || res.City != city || res.Zip != zip || res.County != county ||
		res.Country != country || res.State != state {
		t.Error("Invalid address from db, expected username")
	}

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to create with no valid post form
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)

	req.PostForm = url.Values{}

	address.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestAddress_View(t *testing.T) {
	//Create address with params
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)
	prepareRequest(req)

	address.Create()

	res := &models.Address{}
	ctx.DB.Find(res, "number = ?", number)

	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	address.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to view without params id
	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	address.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	req, rr = prepareReqAndRecorder("GET", "/address/view")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	address.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

func TestAddress_Edit(t *testing.T) {
	//Create address with params
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)
	prepareRequest(req)

	address.Create()

	res := &models.Address{}
	ctx.DB.Find(res, "number = ?", number)

	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)
	prepareRequest(req)

	address.Edit()

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to post data without id
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	address.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid field
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	address.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid data)
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("number", strconv.Itoa(111111))
	req.PostForm.Add("street", "dafdsf")

	address.Edit()

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to update non exist mail
	req, rr = prepareReqAndRecorder("POST", "/address/update")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	address.Edit()

	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

func TestAddress_Delete(t *testing.T) {
	//Create address with params
	req, rr = prepareReqAndRecorder("POST", "/address/create")
	address, ctx = prepareAddress(req, rr)
	prepareRequest(req)

	address.Create()

	res := &models.Address{}
	ctx.DB.Find(res, "number = ?", number)

	req, rr = prepareReqAndRecorder("GET", "/address/delete")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	address.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to delete non exist mail
	req, rr = prepareReqAndRecorder("GET", "/address/delete")
	address, ctx = prepareAddress(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(12345678)

	address.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}

	//Try to delete without param id
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
	conf, err := config.NewConfig(path)

	if err != nil {
		panic(err)
	}

	err = model.OpenWithConfig(conf)

	if err != nil {
		panic(err)
	}

	ctx := base.NewContext(rr, req)
	ctx.DB = model
	ctx.Log = logger.NewDefaultLogger(os.Stdout)

	contr := BaseController{Ctx: ctx, Routes: addressRoutes}
	address := &Address{}

	address.Routes = emailRoutes
	address.BaseController = contr

	return address, ctx
}

func prepareRequest(req *http.Request) {
	req.PostForm = url.Values{}
	req.PostForm.Add("number", strconv.Itoa(number))
	req.PostForm.Add("street", street)
	req.PostForm.Add("state", state)
	req.PostForm.Add("city", city)
	req.PostForm.Add("zip", zip)
	req.PostForm.Add("county", county)
	req.PostForm.Add("country", country)
}
