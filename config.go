package utron

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"reflect"
)

// Config stores configurations values
type Config struct {
	AppName      string `json:"app_name" yaml:"app_name" toml:"app_name"`
	BaseURL      string `json:"base_url" yaml:"base_url" toml:"base_url"`
	Port         int    `json:"port" yaml:"port" toml:"port"`
	Verbose      bool   `json:"verbose" yaml:"verbose" toml:"verbose"`
	StaticDir    string `json:"static_dir" yaml:"static_dir" toml:"static_dir"`
	ViewsDir     string `json:"view_dir" yaml:"view_dir" toml:"view_dir"`
	Database     string `json:"database" yaml:"database" toml:"database"`
	DatabaseConn string `json:"database_conn" yaml:"database_conn" toml:"database_conn"`
}

// DefaultConfig returns the default configuation settings.
func DefaultConfig() *Config {
	return &Config{
		AppName:   "utron web app",
		BaseURL:   "http://localhost:8090",
		Port:      8090,
		Verbose:   false,
		StaticDir: "static",
		ViewsDir:  "views",
	}
}

// NewConfig reads configuration from path. The format is deductted from file extension
//	* .json    - is decoded as json
//	* .yml     - is decoded as yaml
//	* .toml    - is decoded as toml
func NewConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	switch filepath.Ext(path) {
	case ".json":
		jerr := json.Unmarshal(data, cfg)
		if jerr != nil {
			return nil, jerr
		}
	case ".toml":
		_, terr := toml.Decode(string(data), cfg)
		if terr != nil {
			return nil, terr
		}
	case ".yml":
		yerr := yaml.Unmarshal(data, cfg)
		if yerr != nil {
			return nil, yerr
		}

	default:
		return nil, errors.New("utron: config file not supported")
	}



	return cfg, nil
}

// This function overrides some values on the config with environement variables
// The idea is to provide a way to use utron in different environents.
//
// - PORT
// - BASE_URL
// - DATABASE
// - DATABASE_CONN

func (configuration *Config) ApplyEnvironmentVariables() {
	stringVariableNames := map[string]string{
		"BASE_URL": "BaseURL",
		"APP_NAME":	"AppName",
		"VERBOSE":	"Verbose",
		"STATIC_DIR":	"StaticDir",
		"VIEW_DIR":	"ViewsDir",
		"DATABASE": "Database",
		"DATABASE_CONN" : "DatabaseConn",
	}

	environmentValue := os.Getenv("PORT")
	if environmentValue != "" {
		configuration.Port, _ = strconv.Atoi(environmentValue)
	}

	for envVariableName, fieldName := range stringVariableNames {
		environmentValue = os.Getenv(envVariableName)
		if environmentValue != "" {
			configuration.setStringFieldValue(fieldName, environmentValue)
		}
	}
}

func (configuration *Config) setStringFieldValue(fieldName string, value string){
	ps := reflect.ValueOf(configuration)
	s := ps.Elem()

	field := s.FieldByName(fieldName)
	field.SetString(value)
}

// saveToFile saves the Config in the file named path. This is a helper method
// for generating sample configuration files.
func (c *Config) saveToFile(path string) error {
	var data []byte
	switch filepath.Ext(path) {
	case ".json":
		d, err := json.MarshalIndent(c, "", "\t") // use tab indent to make it human friendly
		if err != nil {
			return err
		}
		data = d
	case ".yml":
		d, err := yaml.Marshal(c)
		if err != nil {
			return err
		}
		data = d
	case ".toml":
		b := &bytes.Buffer{}
		err := toml.NewEncoder(b).Encode(c)
		if err != nil {
			return err
		}
		data = b.Bytes()

	}
	return ioutil.WriteFile(path, data, 0600)
}
