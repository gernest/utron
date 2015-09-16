package utron

import "testing"

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
		if nCfg.AppName != nCfg.AppName {
			t.Errorf("expecetd %s got %s", cfg.AppName, nCfg.AppName)
		}
	}

}
