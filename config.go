package utron

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/camelcase"
	"gopkg.in/yaml.v2"
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
	Automigrate  bool   `json:"automigrate" yaml:"automigrate" toml:"automigrate"`
}

// DefaultConfig returns the default configuration settings.
func DefaultConfig() *Config {
	return &Config{
		AppName:     "utron web app",
		BaseURL:     "http://localhost:8090",
		Port:        8090,
		Verbose:     false,
		StaticDir:   "static",
		ViewsDir:    "views",
		Automigrate: true,
	}
}

// NewConfig reads configuration from path. The format is deduced from the file extension
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
		return nil, errors.New("utron: config file format not supported")
	}
	return cfg, nil
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

// SyncEnv overrides c field's values that are set in the environment.
//
// The environment variable names are derived from config fields by underscoring, and uppercasing
// the name. E.g. AppName will have a corresponding environment variable APP_NAME
//
// NOTE only int, string and bool fields are supported and the corresponding values are set.
// when the field value is not supported it is ignored.
func (c *Config) SyncEnv() error {
	cfg := reflect.ValueOf(c).Elem()
	cTyp := cfg.Type()

	for k := range make([]struct{}, cTyp.NumField()) {
		field := cTyp.Field(k)

		cm := getEnvName(field.Name)
		env := os.Getenv(cm)
		if env == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			cfg.FieldByName(field.Name).SetString(env)
		case reflect.Int:
			v, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("utron: loading config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).Set(reflect.ValueOf(v))
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("utron: loading config field %s %v", field.Name, err)
			}
			cfg.FieldByName(field.Name).SetBool(b)
		}

	}
	return nil
}

// getEnvName returns all upper case and underscore separated string, from field.
// field is a camel case string.
//
// example
//	AppName will change to APP_NAME
func getEnvName(field string) string {
	camSplit := camelcase.Split(field)
	var rst string
	for k, v := range camSplit {
		if k == 0 {
			rst = strings.ToUpper(v)
			continue
		}
		rst = rst + "_" + strings.ToUpper(v)
	}
	return rst
}
