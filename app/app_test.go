package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gernest/utron/config"
	"github.com/gernest/utron/controller"
)

const notFoundMsg = "nothing"

func TestGetAbsPath(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// non existing
	_, err = getAbsolutePath("nope")
	if err == nil {
		t.Error("expcted error got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expcetd not exist got %v", err)
	}

	absPath := filepath.Join(wd, "fixtures")

	// Relqtive
	dir, err := getAbsolutePath("fixtures")
	if err != nil {
		t.Error(err)
	}

	if dir != absPath {
		t.Errorf("expceted %s got %s", absPath, dir)
	}

	// Absolute
	dir, err = getAbsolutePath(absPath)
	if err != nil {
		t.Error(err)
	}

	if dir != absPath {
		t.Errorf("expceted %s got %s", absPath, dir)
	}

}

type SimpleMVC struct {
	controller.BaseController
}

func (s *SimpleMVC) Hello() {
	s.Ctx.Data["Name"] = "gernest"
	s.Ctx.Template = "index"
	s.String(http.StatusOK)
}

func TestMVC(t *testing.T) {
	app, err := NewMVC("fixtures/mvc")
	if err != nil {
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

func TestApp(t *testing.T) {
	app := NewApp()
	// Set not found handler
	err := app.SetNotFoundHandler(http.HandlerFunc(sampleDefault))
	if err != nil {
		t.Error(err)
	}

	// no router
	app.Router = nil
	err = app.SetNotFoundHandler(http.HandlerFunc(sampleDefault))
	if err == nil {
		t.Error("expected an error")
	}
}

func sampleDefault(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(notFoundMsg))
}

func TestStaticServer(t *testing.T) {
	c := &config.Config{}
	_, ok, _ := StaticServer(c)
	if ok {
		t.Error("expected false")
	}
	c.StaticDir = "fixtures"
	s, ok, _ := StaticServer(c)
	if !ok {
		t.Error("expected true")
	}
	expect := "/static/"
	if s != expect {
		t.Errorf("expected %s got %s", expect, s)
	}
}
