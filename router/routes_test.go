package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gernest/utron/controller"
)

var msg = "gernest"

type Sample struct {
	controller.BaseController
	Routes []string
}

func (s *Sample) Bang() {
	_, _ = s.Ctx.Write([]byte(msg))
	s.JSON(http.StatusOK)
}

func (s *Sample) Hello() {
	_, _ = s.Ctx.Write([]byte(msg))
	s.String(http.StatusOK)
}

func NewSample() *Sample {
	routes := []string{
		"get,post;/hello/world;Hello",
	}
	s := &Sample{}
	s.Routes = routes
	return s
}

func TestRouterAdd(t *testing.T) {
	r := NewRouter()
	_ = r.Add(controller.GetCtrlFunc(&Sample{}))

	req, err := http.NewRequest("GET", "/sample/bang", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}
}

func TestRouteField(t *testing.T) {
	r := NewRouter()
	routes := []string{
		"get,post;/hello/world;Hello",
	}
	s := &Sample{}
	s.Routes = routes
	err := r.Add(controller.GetCtrlFunc(s))
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/hello/world", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}

	req, err = http.NewRequest("GET", "/sample/bang", nil)
	if err != nil {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected %s got %s", msg, w.Body.String())
	}
}

func TestRoutesFile(t *testing.T) {

	routeFiles := []string{
		"../fixtures/config/routes.json",
		"../fixtures/config/routes.yml",
		"../fixtures/config/routes.toml",
		"../fixtures/config/routes.hcl",
	}

	for _, file := range routeFiles {
		r := NewRouter()

		err := r.LoadRoutesFile(file)
		if err != nil {
			t.Error(err)
		}
		if len(r.routes) != 2 {
			t.Errorf("expcted 2 got %d", len(r.routes))
		}
		_ = r.Add(controller.GetCtrlFunc(NewSample()))

		req, err := http.NewRequest("GET", "/hello", nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected %d got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != msg {
			t.Errorf("expected %s got %s", msg, w.Body.String())
		}
	}

}

func TestSplitRoutes(t *testing.T) {
	data := []struct {
		routeStr, ctrl, fn string
	}{
		{
			"get,post;/;Hello.Home", "Hello", "Home",
		},
		{
			"get,post;/;Home", "", "Home",
		},
	}

	for _, v := range data {
		r, err := splitRoutes(v.routeStr)
		if err != nil {
			t.Fatal(err)
		}
		if r.ctrl != v.ctrl {
			t.Errorf("expected %s got %s", v.ctrl, r.ctrl)
		}
		if r.fn != v.fn {
			t.Errorf("extected %s got %s", v.fn, r.fn)
		}
	}
}
