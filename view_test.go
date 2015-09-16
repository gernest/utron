package utron

import (
	"bytes"
	"strings"
	"testing"
)

func TestSimpleView(t *testing.T) {
	v, err := NewSimpleView("fixtures/view")
	if err != nil {
		t.Error(err)
	}

	out := &bytes.Buffer{}
	data := struct {
		Name string
	}{
		"gernest",
	}

	tpls := []string{
		"index", "sample/hello",
	}
	for _, tpl := range tpls {
		verr := v.Render(out, tpl, data)
		if verr != nil {
			t.Error(err)
		}
		if !strings.Contains(out.String(), data.Name) {
			t.Errorf("expeted %s to contain %s", out.String(), data.Name)
		}
	}

}
