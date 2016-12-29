package flash

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gernest/utron/controller"
	"github.com/gernest/utron/logger"
	"github.com/gernest/utron/router"
	"github.com/gorilla/sessions"
)

type FlashTest struct {
	controller.BaseController
	Routes []string
}

const (
	fname = "flash"
	fkey  = "flash"
)

var result Flashes

func (f *FlashTest) Index() {
	fl := New()
	fl.Success("Success")
	fl.Err("Err")
	fl.Warn("Warn")
	fl.Save(f.Ctx, fname, fkey)
}

func (f FlashTest) Flash() {
	r, err := GetFlashes(f.Ctx, fname, fkey)
	if err != nil {
		f.Ctx.Log.Errors(err)
		return
	}
	result = r
}

func NewFlashTest() controller.Controller {
	return &FlashTest{
		Routes: []string{
			"get;/;Index",
			"get;/flash;Flash",
		},
	}
}

func TestFlash(t *testing.T) {
	codecKey1 := "ePAPW9vJv7gHoftvQTyNj5VkWB52mlza"
	codecKey2 := "N8SmpJ00aSpepNrKoyYxmAJhwVuKEWZD"
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{Jar: jar}
	o := &router.Options{
		Log:          logger.NewDefaultLogger(os.Stdout),
		SessionStore: sessions.NewCookieStore([]byte(codecKey1), []byte(codecKey2)),
	}
	r := router.NewRouter(o)
	r.Add(NewFlashTest)
	ts := httptest.NewServer(r)
	defer ts.Close()
	_, err = client.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Error(err)
	}
	_, err = client.Get(fmt.Sprintf("%s/flash", ts.URL))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 got %d", len(result))
	}
}
