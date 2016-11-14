package utron

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gernest/utron/controller"
)

type SimpleMVC struct {
	*controller.BaseController
}

func (s *SimpleMVC) Hello() {
	s.Ctx.Data["Name"] = "gernest"
	s.Ctx.Template = "index"
	s.String(http.StatusOK)
}

func TestMVC(t *testing.T) {
	app, err := NewMVC("fixtures/mvc")
	if err != nil {
		fmt.Println(err)
		t.Skip(err)
	}
	app.AddController(controller.GetCtrlFunc(&SimpleMVC{}))

	req, _ := http.NewRequest("GET", "/simplemvc/hello", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expcted %d got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "gernest") {
		t.Errorf("expected %s to contain gernest", w.Body.String())
	}
}
