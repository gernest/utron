package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"regexp"
	"fmt"

	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/config"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	req          *http.Request
	rr           *httptest.ResponseRecorder
	email        *Email
	ctx          *base.Context
	err          error
	mock         sqlmock.Sqlmock
	id          = -273
)
func TestBaseController(t *testing.T) {
	req, _ = http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx = base.NewContext(w, req)

	ctrl := BaseController{}

	if ctrl.New(ctx); ctrl.Ctx == nil {
		t.Error("expected Ctx to be set")
	}

	// HTML response
	ctrl.HTML(http.StatusOK)
	cTyp := w.Header().Get(base.Content.Type)
	if cTyp != base.Content.TextHTML {
		t.Errorf("expecetd %s got %s", base.Content.TextHTML, cTyp)
	}

	// JSON response
	ctrl.JSON(http.StatusOK)
	cTyp = w.Header().Get(base.Content.Type)
	if cTyp != base.Content.Application.JSON {
		t.Errorf("expected %s got %s", base.Content.Application.JSON, cTyp)
	}

	//RenderJSON test
	//TODO: This aint right- fix it
	ctrl.RenderJSON(*config.DefaultConfig(), http.StatusOK)
	if cTyp != base.Content.Application.JSON {
		t.Errorf("expected %s got %s", base.Content.Application.JSON, cTyp)
	}

	// Plain text response
	ctrl.String(http.StatusOK)
	cTyp = w.Header().Get(base.Content.Type)
	if cTyp != base.Content.TextPlain {
		t.Errorf("expected %s got %s", base.Content.TextPlain, cTyp)
	}

	err = ctrl.Render()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	err = ctrl.Render()
	if err == nil {
		t.Error("expected error got nil")
	}

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

func fixedFullRe(s string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(s))
}