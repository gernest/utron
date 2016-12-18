package base

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gernest/utron/config"
	"github.com/gorilla/mux"
)

type DummyView struct {
}

func (d *DummyView) Render(out io.Writer, name string, data interface{}) error {
	out.Write([]byte(name))
	return nil
}

func TestContext(t *testing.T) {
	r := mux.NewRouter()
	name := "world"
	r.HandleFunc("/hello/{name}", testHandler(t, name))
	req, _ := http.NewRequest("GET", "/hello/"+name, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func testHandler(t *testing.T, name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxHandler(t, name, w, r)
	}
}

func ctxHandler(t *testing.T, name string, w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	ctx.Init()
	pname := ctx.Params["name"]
	if pname != name {
		t.Error("expected %s got %s", name, pname)
	}

	ctx.SetData("name", pname)

	data := ctx.GetData("name")
	if data == nil {
		t.Error("expected values to be stored in context")
	}
	ctx.JSON()
	h := w.Header().Get(Content.Type)
	if h != Content.Application.JSON {
		t.Errorf("expected %s got %s", Content.Application.JSON, h)
	}
	ctx.HTML()
	h = w.Header().Get(Content.Type)
	if h != Content.TextHTML {
		t.Errorf("expected %s got %s", Content.TextHTML, h)
	}
	ctx.TextPlain()
	h = w.Header().Get(Content.Type)
	if h != Content.TextPlain {
		t.Errorf("expected %s got %s", Content.TextPlain, h)
	}

	err := ctx.Commit()
	if err != nil {
		t.Error(err)
	}

	// make sure we can't commit twice
	err = ctx.Commit()
	if err == nil {
		t.Error("expected error")
	}

	// when there is template and view
	ctx.isCommited = false
	ctx.Template = pname
	ctx.Set(&DummyView{})
	ctx.Cfg = &config.Config{}
	err = ctx.Commit()
	if err != nil {
		t.Error(err)
	}
}
