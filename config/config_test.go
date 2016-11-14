package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	cfgFiles := []string{
		"../fixtures/config/app.json",
		"../fixtures/config/app.yml",
		"../fixtures/config/app.toml",
		"../fixtures/config/app.hcl",
	}
	badCfgFiles := []string{
		"../fixtures/badconfig/app.json",
		"../fixtures/badconfig/app.yml",
		"../fixtures/badconfig/app.toml",
		"../fixtures/badconfig/app.hcl",
	}
	for _, f := range badCfgFiles {
		_, err := NewConfig(f)
		if err == nil {
			t.Fatal("expected error ", f)
		}

	}

	cfg := DefaultConfig()

	for _, f := range cfgFiles {
		nCfg, err := NewConfig(f)
		if err != nil {
			t.Fatal(err)
		}
		if nCfg.AppName != cfg.AppName {
			t.Errorf("expecetd %s got %s", cfg.AppName, nCfg.AppName)
		}
	}

	// no file
	_, err := NewConfig("nothing")
	if err == nil {
		t.Error("expected error")
	}

	//unsupporte file
	_, err = NewConfig("../fixtures/view/index.tpl")
	if err == nil {
		t.Error("expected error")
	}
	if err != errCfgUnsupported {
		t.Errorf("expected %v got %v", errCfgUnsupported, err)
	}

}

func TestConfigEnv(t *testing.T) {
	fields := []struct {
		name, env, value string
	}{
		{"AppName", "APP_NAME", "utron"},
		{"BaseURL", "BASE_URL", "http://localhost:8090"},
		{"Port", "PORT", "8091"},
		{"ViewsDir", "VIEWS_DIR", "fixtures/view"},
		{"StaticDir", "STATIC_DIR", "fixtures/todo/static"},
		{"Database", "DATABASE", "postgres"},
		{"DatabaseConn", "DATABASE_CONN", "postgres://postgres@localhost/utron?sslmode=disable"},
		{"Automigrate", "AUTOMIGRATE", "true"},
	}
	for _, f := range fields {

		// check out env name maker
		cm := getEnvName(f.name)
		if cm != f.env {
			t.Errorf("expected %s got %s", f.env, cm)
		}
	}

	// set environment values
	for _, f := range fields {
		_ = os.Setenv(f.env, f.value)
	}

	cfg := DefaultConfig()
	if err := cfg.SyncEnv(); err != nil {
		t.Errorf("syncing env %v", err)
	}

	if cfg.Port != 8091 {
		t.Errorf("expected 8091 got %d instead", cfg.Port)
	}

	if cfg.AppName != "utron" {
		t.Errorf("expected utron got %s", cfg.AppName)
	}
}
