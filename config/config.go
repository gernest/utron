package config

import (
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
	"github.com/gorilla/securecookie"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

var errCfgUnsupported = errors.New("utron: config file format not supported")

// Config stores configurations values
type Config struct {
	AppName      string `json:"app_name" yaml:"app_name" toml:"app_name" hcl:"app_name"`
	BaseURL      string `json:"base_url" yaml:"base_url" toml:"base_url" hcl:"base_url"`
	Port         int    `json:"port" yaml:"port" toml:"port" hcl:"port"`
	Verbose      bool   `json:"verbose" yaml:"verbose" toml:"verbose" hcl:"verbose"`
	StaticDir    string `json:"static_dir" yaml:"static_dir" toml:"static_dir" hcl:"static_dir"`
	ViewsDir     string `json:"view_dir" yaml:"view_dir" toml:"view_dir" hcl:"view_dir"`
	Database     string `json:"database" yaml:"database" toml:"database" hcl:"database"`
	DatabaseConn string `json:"database_conn" yaml:"database_conn" toml:"database_conn" hcl:"database_conn"`
	Automigrate  bool   `json:"automigrate" yaml:"automigrate" toml:"automigrate" hcl:"automigrate"`
	NoModel      bool   `json:"no_model" yaml:"no_model" toml:"no_model" hcl:"no_model"`

	// session
	SessionName     string `json:"session_name" yaml:"session_name" toml:"session_name" hcl:"session_name"`
	SessionPath     string `json:"session_path" yaml:"session_path" toml:"session_path" hcl:"session_path"`
	SessionDomain   string `json:"session_domain" yaml:"session_domain" toml:"session_domain" hcl:"session_domain"`
	SessionMaxAge   int    `json:"session_max_age" yaml:"session_max_age" toml:"session_max_age" hcl:"session_max_age"`
	SessionSecure   bool   `json:"session_secure" yaml:"session_secure" toml:"session_secure" hcl:"session_secure"`
	SessionHTTPOnly bool   `json:"session_httponly" yaml:"session_httponly" toml:"session_httponly" hcl:"session_httponly"`

	// The name of the session store to use
	// Options are
	// file , cookie ,ql
	SessionStore string `json:"session_store" yaml:"session_store" toml:"session_store" hcl:"session_store"`

	// Flash is the session name for flash messages
	Flash string `json:"flash" yaml:"flash" toml:"flash" hcl:"flash"`

	// KeyPair for secure cookie its a comma separates strings of keys.
	SessionKeyPair []string `json:"session_key_pair" yaml:"session_key_pair" toml:"session_key_pair" hcl:"session_key_pair"`

	// flash message
	FlashContextKey string `json:"flash_context_key" yaml:"flash_context_key" toml:"flash_context_key" hcl:"flash_context_key"`
}

// DefaultConfig returns the default configuration settings.
func DefaultConfig() *Config {
	a := securecookie.GenerateRandomKey(32)
	b := securecookie.GenerateRandomKey(32)
	return &Config{
		AppName:       "utron web app",
		BaseURL:       "http://localhost:8090",
		Port:          8090,
		Verbose:       false,
		StaticDir:     "static",
		ViewsDir:      "views",
		Automigrate:   true,
		SessionName:   "_utron",
		SessionPath:   "/",
		SessionMaxAge: 2592000,
		SessionKeyPair: []string{
			string(a), string(b),
		},
		Flash: "_flash",
	}
}

// NewConfig reads configuration from path. The format is deduced from the file extension
//	* .json    - is decoded as json
//	* .yml     - is decoded as yaml
//	* .toml    - is decoded as toml
//  * .hcl	   - is decoded as hcl
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
	case ".hcl":
		obj, herr := hcl.Parse(string(data))
		if herr != nil {
			return nil, herr
		}
		if herr = hcl.DecodeObject(&cfg, obj); herr != nil {
			return nil, herr
		}
	default:
		return nil, errCfgUnsupported
	}
	err = cfg.SyncEnv()
	if err != nil {
		return nil, err
	}

	// ensure the key pairs are set
	if cfg.SessionKeyPair == nil {
		a := securecookie.GenerateRandomKey(32)
		b := securecookie.GenerateRandomKey(32)
		cfg.SessionKeyPair = []string{
			string(a), string(b),
		}
	}
	return cfg, nil
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
