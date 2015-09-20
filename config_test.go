package utron

import (
	"testing"
	"os"
)

func TestConfig(t *testing.T) {
	cfgFiles := []string{
		"fixtures/config/app.json",
		"fixtures/config/app.yml",
		"fixtures/config/app.toml",
	}

	cfg := DefaultConfig()

	//	//	 Uncomment the this to generate the sample
	//	//	  config files in fixures/config directory
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

func TestEnvironmentOverride(t *testing.T){
	err := os.Setenv("PORT", "8888")
	err = os.Setenv("DATABASE_CONN", "myDatabase://hey:cool@databse.com/my_database")
	err = os.Setenv("DATABASE", "myDatabase")

	if err == nil {

		cfg := DefaultConfig()
		cfg,  err = NewConfig("fixtures/config/app.json")
		cfg.SyncEnv()

		if cfg.Port != 8888 {
			t.Errorf("expecetd %s got %s", 8888, cfg.Port)
		}

		if cfg.Database != "myDatabase" {
			t.Errorf("expecetd %s got %s", "myDatabase", cfg.Database)
		}
	}
}

func TestEnvironmentOverrideWithInvalidValues(t *testing.T){
	err := os.Setenv("PORT", "aaa")
	err = os.Setenv("VERBOSE", "NON-BOOL")

	if err == nil {

		cfg := DefaultConfig()
		cfg,  err = NewConfig("fixtures/config/app.json")
		cfg.SyncEnv()

		if cfg.Port != 8090 {
			t.Errorf("expecetd %s got %s", 8090, cfg.Port)
		}

		if cfg.Verbose != false {
			t.Errorf("expecetd %s got %s", false, cfg.Verbose)
		}
	}
}
