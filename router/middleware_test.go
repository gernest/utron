package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gernest/utron/base"
	"github.com/gernest/utron/controller"
	"github.com/gorilla/context"
)

const incrementKey = "increment"

func (s *Sample) Increment() {
	key := s.Ctx.GetData(incrementKey)
	fmt.Fprintf(s.Ctx, "%v", key)
	s.String(http.StatusOK)
}

func plainIncrement(n int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key, ok := context.GetOk(r, incrementKey); ok {
				ikey := key.(int)
				ikey += n
				context.Set(r, incrementKey, ikey)
			} else {
				context.Set(r, incrementKey, 0)
			}
			h.ServeHTTP(w, r)
		})
	}
}

func contextMiddleware(n int) func(*base.Context) error {
	return func(ctx *base.Context) error {
		key := ctx.GetData(incrementKey)
		if key != nil {
			ikey := key.(int)
			ikey += n
			ctx.SetData(incrementKey, ikey)
		} else {
			ctx.SetData(incrementKey, 0)
		}
		return nil
	}
}

func TestMiddlewarePlain(t *testing.T) {
	expect := "3"
	r := NewRouter()
	_ = r.Add(controller.GetCtrlFunc(&Sample{}), plainIncrement(0), plainIncrement(1), plainIncrement(2))

	req, err := http.NewRequest("GET", "/sample/increment", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != expect {
		t.Errorf("expected %s got %s", expect, w.Body.String())
	}
}

func TestMiddlewareContext(t *testing.T) {
	expect := "3"
	r := NewRouter()
	_ = r.Add(controller.GetCtrlFunc(&Sample{}), contextMiddleware(0), contextMiddleware(1), contextMiddleware(2))

	req, err := http.NewRequest("GET", "/sample/increment", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != expect {
		t.Errorf("expected %s got %s", expect, w.Body.String())
	}
}

func TestMiddlewareMixed(t *testing.T) {
	expect := "6"

	r := NewRouter()
	_ = r.Add(controller.GetCtrlFunc(&Sample{}), plainIncrement(0), contextMiddleware(1), plainIncrement(2), contextMiddleware(3))

	req, err := http.NewRequest("GET", "/sample/increment", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != expect {
		t.Errorf("expected %s got %s", expect, w.Body.String())
	}
}
