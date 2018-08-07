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

const path = "../fixtures/test/app.json"

var (
	emailRoutes = []string{
		"get;/email;Index",
		"post;/email/create;Create",
		"get;/email/view/{id};View",
		"get;/email/delete/{id};Delete",
		"get;/email/update/{id};ViewEdit",
		"post;/email/update/{id};Edit",
	}
	username = "someadress"
	domain   = "gmail.com"
	req      *http.Request
	rr       *httptest.ResponseRecorder
	email    *Email
	ctx      *base.Context
)

func TestEmail_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/email")
	email, _ = prepareEmail(req, rr)

	email.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestEmail_Create(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	email.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Create mail with address: someadress@gmail.com
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+"@"+domain)

	email.Create()

	res := &models.Email{}
	ctx.DB.Find(res, "address = ?", username+"@"+domain)

	if res == nil || res.Username != username || res.Domain != domain {
		t.Error("Invalid email from db, expected username: ", username, " domain: ", domain, " got: ",
			res.Username, " ", res.Domain)
	}

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to create with no valid post form
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}

	email.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestEmail_View(t *testing.T) {
	//Create mail with address: someadress@gmail.com
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+"@"+domain)

	email.Create()

	res := &models.Email{}
	ctx.DB.Find(res, "address = ?", username+"@"+domain)

	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	email.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to view without params id
	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareEmail(req, rr)

	email.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	req, rr = prepareReqAndRecorder("GET", "/email/view")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	email.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

//func TestEmail_ViewEdit(t *testing.T) {
//	//Create mail with address: someadress@gmail.com
//	req, rr = prepareReqAndRecorder("POST", "/email/create")
//	email, ctx = prepareEmail(req, rr)
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
//	email, ctx = prepareEmail(req, rr)
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
//	email, ctx = prepareEmail(req, rr)
//
//	email.ViewEdit()
//	if rr.Result().StatusCode != http.StatusInternalServerError {
//		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
//	}
//
//	//Try to viewEdit non exist email
//	req, rr = prepareReqAndRecorder("GET", "/email/update")
//	email, ctx = prepareEmail(req, rr)
//
//	ctx.Params = make(map[string]string)
//	ctx.Params["id"] = strconv.Itoa(123456784)
//
//	email.ViewEdit()
//	if rr.Result().StatusCode != http.StatusNotFound {
//		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
//	}
//}

func TestEmail_Edit(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+"@"+domain)

	email.Create()

	res := &models.Email{}
	ctx.DB.Find(res, "address = ?", username+"@"+domain)

	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)
	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+"sdasd"+"@"+domain)

	email.Edit()

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to post data without id
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareEmail(req, rr)

	email.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid field
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	email.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid data(email)
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+domain)

	email.Edit()

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to update non exist mail
	req, rr = prepareReqAndRecorder("POST", "/email/update")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	email.Edit()

	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

func TestEmail_Delete(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", "/email/create")
	email, ctx = prepareEmail(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("address", username+"@"+domain)

	email.Create()

	res := &models.Email{}
	ctx.DB.Find(res, "address = ?", username+"@"+domain)

	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	email.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to delete non exist mail
	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareEmail(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(12345678)

	email.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}

	//Try to delete without param id
	req, rr = prepareReqAndRecorder("GET", "/email/delete")
	email, ctx = prepareEmail(req, rr)

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
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()

	return req, rr
}

//Create *Email with BaseConroller(with DB) and emailRoutes
func prepareEmail(req *http.Request, rr *httptest.ResponseRecorder) (*Email, *base.Context) {
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

	contr := BaseController{Ctx: ctx, Routes: emailRoutes}
	email := &Email{}

	email.Routes = emailRoutes
	email.BaseController = contr

	return email, ctx
}
