package utron

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseController(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req)

	ctrl := &BaseController{}

	if ctrl.New(ctx); ctrl.Ctx == nil {
		t.Error("expected Ctx to be set")
	}

	// HTML response
	ctrl.HTML(http.StatusOK)
	cTyp := w.Header().Get(Content.Type)
	if cTyp != Content.TextHTML {
		t.Errorf("expecetd %s got %s", Content.TextHTML, cTyp)
	}

	// JSON response
	ctrl.JSON(http.StatusOK)
	cTyp = w.Header().Get(Content.Type)
	if cTyp != Content.Application.JSON {
		t.Errorf("expected %s got %s", Content.Application.JSON, cTyp)
	}

	// Plain text response
	ctrl.String(http.StatusOK)
	cTyp = w.Header().Get(Content.Type)
	if cTyp != Content.TextPlain {
		t.Errorf("expected %s got %s", Content.TextPlain, cTyp)
	}

	err := ctrl.Render()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	err = ctrl.Render()
	if err == nil {
		t.Error("expected error got nil")
	}

}
