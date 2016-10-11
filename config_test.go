package utron

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	cfgFiles := []string{
		"fixtures/config/app.json",
		"fixtures/config/app.yml",
		"fixtures/config/app.toml",
		"fixtures/config/app.hcl",
	}

	cfg := DefaultConfig()

	//	//	 Uncomment this to generate the sample
	//	//	 config files in the fixures/config directory
	//	cfg.ViewsDir = "fixtures/view"
	//	cfg.StaticDir = "fixtures/static"
	//	for _, f := range cfgFiles {
	//		err := cfg.saveToFile(f)
	//		if err != nil {
	//			t.Error(err)
	//		}
	//	}

	for _, f := range cfgFiles {
		nCfg, err := NewConfig(f)
		if err != nil {
			t.Fatal(err)
		}
		if nCfg.AppName != cfg.AppName {
			t.Errorf("expecetd %s got %s", cfg.AppName, nCfg.AppName)
		}
	}

}

func TestConfigEnv(t *testing.T) {
	fields := []struct {
		name, env, value string
	}{
		{"AppName", "APP_NAME", "utron"},
		{"BaseURL", "BASE_URL", "http://localhost:8090"},
		{"Port", "PORT", "8091"},
		{"ViewsDir", "VIEWS_DIR", "viewTest"},
		{"StaticDir", "STATIC_DIR", "statics"},
		{"Database", "DATABASE", "utro_db"},
		{"DatabaseConn", "DATABASE_CONN", "mydb_conn"},
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
