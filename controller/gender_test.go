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
	genderRoutes = []string{
		"get;/gender;Index",
		"post;/gender/create;Create",
		"get;/gender/view/{id};View",
		"get;/gender/delete/{id};Delete",
	}
	gender     *Gender
	claimedSex = "whatisit"
	bioSex     = byte(0)
)

func TestGender_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", "/gender")
	gender, _ = prepareGender(req, rr)

	gender.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestGender_Create(t *testing.T) {
	//Try to create gender with invalid form params
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	gender.Create()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Create gender with params
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)
	prepareRequestGender(req)

	gender.Create()

	res := &models.Gender{}
	ctx.DB.Find(res, "bio_sex = ?", bioSex)

	if res.ClaimedSex != claimedSex || res.BioSex != bioSex {
		t.Error("Invalid gender from db, expected username")
	}

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to create with no valid post form
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)

	req.PostForm = url.Values{}

	gender.Create()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestGender_View(t *testing.T) {
	//Create gender with params
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)
	prepareRequestGender(req)

	gender.Create()

	res := &models.Gender{}
	ctx.DB.Find(res, "bio_sex = ?", bioSex)

	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	gender.View()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to view without params id
	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	gender.View()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	req, rr = prepareReqAndRecorder("GET", "/gender/view")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	gender.View()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

func TestGender_Edit(t *testing.T) {
	//Create gender with params
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)
	prepareRequestGender(req)

	gender.Create()

	res := &models.Gender{}
	ctx.DB.Find(res, "bio_sex = ?", bioSex)

	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)
	prepareRequestGender(req)

	gender.Edit()

	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to post data without id
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	gender.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid field
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("Sdsd", "dsads")

	gender.Edit()

	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}

	//Try to post form with invalid data)
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	req.PostForm = url.Values{}
	req.PostForm.Add("claimed_sex", "")

	gender.Edit()

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}

	ctx.DB.Delete(res)
	//Try to update non exist mail
	req, rr = prepareReqAndRecorder("POST", "/gender/update")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(123456784)

	gender.Edit()

	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}
}

func TestGender_Delete(t *testing.T) {
	//Create gender with params
	req, rr = prepareReqAndRecorder("POST", "/gender/create")
	gender, ctx = prepareGender(req, rr)
	prepareRequestGender(req)

	gender.Create()

	res := &models.Gender{}
	ctx.DB.Find(res, "bio_sex = ?", bioSex)

	req, rr = prepareReqAndRecorder("GET", "/gender/delete")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(res.ID)

	gender.Delete()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}

	//Try to delete non exist mail
	req, rr = prepareReqAndRecorder("GET", "/gender/delete")
	gender, ctx = prepareGender(req, rr)

	ctx.Params = make(map[string]string)
	ctx.Params["id"] = strconv.Itoa(12345678)

	gender.Delete()
	if rr.Result().StatusCode != http.StatusNotFound {
		t.Error("Expected http.StatusNotFound, got: ", rr.Result().Status)
	}

	//Try to delete without param id
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

	contr := BaseController{Ctx: ctx, Routes: genderRoutes}
	gender := &Gender{}

	gender.Routes = emailRoutes
	gender.BaseController = contr

	return gender, ctx
}

func prepareRequestGender(req *http.Request) {
	req.PostForm = url.Values{}
	req.PostForm.Add("claimed_sex", claimedSex)
	req.PostForm.Add("legal_sex", "0")
}
