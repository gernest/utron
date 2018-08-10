# config
--
    import "github.com/NlaakStudios/gowaf/config"


## Usage

#### type Config

```go
type Config struct {
	AppName      string `json:"app_name" yaml:"app_name" toml:"app_name" hcl:"app_name"`
	Domain       string `json:"domain" yaml:"domain" toml:"domain" hcl:"domain"`
	CompanyName  string `json:"company_name" yaml:"company_name" toml:"company_name" hcl:"company_name"`
	BaseURL      string `json:"base_url" yaml:"base_url" toml:"base_url" hcl:"base_url"`
	Port         int    `json:"port" yaml:"port" toml:"port" hcl:"port"`
	Verbose      bool   `json:"verbose" yaml:"verbose" toml:"verbose" hcl:"verbose"`
	FixturesDir  string `json:"fixtures_dir" yaml:"fixtures_dir" toml:"fixtures_dir" hcl:"fixtures_dir"`
	StaticDir    string `json:"static_dir" yaml:"static_dir" toml:"static_dir" hcl:"static_dir"`
	ViewsDir     string `json:"view_dir" yaml:"view_dir" toml:"view_dir" hcl:"view_dir"`
	Database     string `json:"database" yaml:"database" toml:"database" hcl:"database"`
	DatabaseConn string `json:"database_conn" yaml:"database_conn" toml:"database_conn" hcl:"database_conn"`
	Automigrate  bool   `json:"automigrate" yaml:"automigrate" toml:"automigrate" hcl:"automigrate"`
	LoadTestData bool   `json:"load_test_data" yaml:"load_test_data" toml:"load_test_data" hcl:"load_test_data"`
	NoModel      bool   `json:"no_model" yaml:"no_model" toml:"no_model" hcl:"no_model"`
	GoogleID     string `json:"googleid" yaml:"googleid" toml:"googleid" hcl:"googleid"`

	Notifications bool   `json:"notifications" yaml:"notifications" toml:"notifications" hcl:"notifications"`
	Mail          bool   `json:"mail" yaml:"mail" toml:"mail" hcl:"mail"`
	Profile       bool   `json:"profile" yaml:"profile" toml:"profile" hcl:"profile"`
	ThemeColor    string `json:"themecolor" yaml:"themecolor" toml:"themecolor" hcl:"themecolor"`

	FlashTime  uint `json:"flash_time" yaml:"flash_time" toml:"flash_time" hcl:"flash_time"`
	FlashStack uint `json:"flash_stack" yaml:"flash_stack" toml:"flash_stack" hcl:"flash_stack"`
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
```

Config stores configurations values

#### func  DefaultConfig

```go
func DefaultConfig() *Config
```
DefaultConfig returns the default configuration settings.

#### func  NewConfig

```go
func NewConfig(path string) (*Config, error)
```
NewConfig reads configuration from path. The format is deduced from the file
extension

    	* .json    - is decoded as json
    	* .yml     - is decoded as yaml
    	* .toml    - is decoded as toml
     * .hcl	   - is decoded as hcl

#### func (*Config) SyncEnv

```go
func (c *Config) SyncEnv() error
```
SyncEnv overrides c field's values that are set in the environment.

The environment variable names are derived from config fields by underscoring,
and uppercasing the name. E.g. AppName will have a corresponding environment
variable APP_NAME

NOTE only int, string and bool fields are supported and the corresponding values
are set. when the field value is not supported it is ignored.
