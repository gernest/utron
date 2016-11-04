# configuration
utron support yaml, json and toml configurations files. In our todo app, we put the configuration files in the config directory. I have included all three formats for clarity, you can be just fine with either one of them.

`utron` searches for a file named `app.json`, or `app.yml`, `app.toml`, `app.hcl` in the config directory. The first to be found is the one to be used.

This is the content of `config/app.json` file:

```json
{
	"app_name": "utron web app",
	"base_url": "http://localhost:8090",
	"port": 8090,
	"verbose": false,
	"static_dir": "static",
	"view_dir": "views",
	"database": "postgres",
	"database_conn": "postgres://postgres:postgres@localhost/todo",
	"automigrate": true
}
```

You can override the values from the config file by setting environment variables. The names of the environment variables are shown below (with their details)

setting       | environment name | details
--------------|------------------|----------------
app_name      | APP_NAME         | application name
base_url      | BASE_URL         | the base url to use in your views
port          | PORT             | port number the server will listen on
verbose       | VERBOSE          | if set to true, will make all state information log to stdout
static_dir    | STATIC_DIR       | directory to serve static files e.g. images, js or css
view_dir      | VIEWS_DIR        | directory to look for views
database      | DATABASE         | the name of the database you use, e.g. postgres, mysql, sqlite3, foundation
database_conn | DATABASE_CONN    | connection string to your database
automigrate   | AUTOMIGRATE      | creates the tables for models automatically.

If you haven't specified explicitly the location of the configuration directory, it defaults to the directory named `config` in the current working directory.
