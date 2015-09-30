package utron

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type SimpleMVC struct {
	*BaseController
}

func (s *SimpleMVC) Hello() {
	s.Ctx.Data["Name"] = "gernest"
	s.Ctx.Template = "index"
	s.String(http.StatusOK)
}

func TestMVC(t *testing.T) {
	app, err := NewMVC("fixtures/config")
	if err != nil {
		t.Skip(err)
	}
	app.AddController(&SimpleMVC{})

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

func TestGlobal(t *testing.T) {

	//NOTE this is not a serious test
	cfgPath := "fixtures/config"
	SetConfigPath(cfgPath)
	RegisterModels(&Sample{})
	RegisterController(&Sample{})
}
