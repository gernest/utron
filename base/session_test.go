package base

import "testing"

func TestContextSession(t *testing.T) {
	// should error when the session store is not set
	name := "sess"
	ctx := &Context{}
	_, err := ctx.NewSession(name)
	if err == nil {
		t.Error("expected error ", errNoStore)
	}
	_, err = ctx.GetSession(name)
	if err == nil {
		t.Error("expected error ", errNoStore)
	}
	err = ctx.SaveSession(nil)
	if err == nil {
		t.Error("expected error ", errNoStore)
	}
}
