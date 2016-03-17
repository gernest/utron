package utron

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	cfgFiles := []string{
		"fixtures/config/app.json",
		"fixtures/config/app.yml",
		"fixtures/config/app.toml",
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
		os.Setenv(f.env, f.value)
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

func TestConfigEnvEmbeds(t *testing.T) {
	type teststr struct {
		Field1 struct {
			FieldInc string
			Field2   struct {
				FieldInc2 bool
			}
		}
	}

	str := teststr{}
	fields := []struct {
		env, value string
		expected   interface{}
	}{
		{"FIELD_1_FIELD_INC", "fieldv1", "fieldv1"},
		{"FIELD_1_FIELD_2_FIELD_INC_2", "true", true},
	}

	// set envvars
	for _, f := range fields {
		os.Setenv(f.env, f.value)
	}

	cfg := reflect.ValueOf(&str).Elem()
	if err := syncEnv(cfg, ""); err != nil {
		t.Errorf("syncing embedded env %v", err)
	}

	if str.Field1.FieldInc != fields[0].expected {
		t.Errorf("expected %s got %s", fields[0].value, str.Field1.FieldInc)
	}

	if str.Field1.Field2.FieldInc2 != fields[1].expected {
		t.Errorf("expected %s got %s", fields[1].value, str.Field1.Field2.FieldInc2)
	}
}
