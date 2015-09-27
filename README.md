# utron [![GoDoc](https://godoc.org/github.com/gernest/utron?status.svg)](https://godoc.org/github.com/gernest/utron) [![Coverage Status](https://coveralls.io/repos/gernest/utron/badge.svg?branch=master&service=github)](https://coveralls.io/github/gernest/utron?branch=master) [![Build Status](https://travis-ci.org/gernest/utron.svg)](https://travis-ci.org/gernest/utron)
`utron` is a lightweight MVC framework in Go ([Golang](https://golang.org)) for building fast, scalable and robust database-driven web applications.

# Features
* [x] Postgres, MySQL and Foundation database support
* [x] Modular (you can choose which components to use)
* [x] Middleware support. All [alice](https://github.com/justinas/alice) compatible Middleware works out of the box
* [x] Gopher spirit (write golang, use all the golang libraries you like)
* [x] Lightweight. Only MVC
* [x] Multiple configuration files support (currently json, yaml and toml)



# Overview
`utron` is a lightweight MVC framework. It is based on the principles of simplicity, relevance and elegance.

* Simplicity. The design is simple, easy to understand and doesn't introduce many layers between you and the standard library. It is a goal of the project that users should be able to understand the whole framework in a single day.

* Relevance. `utron` doesn't assume anything. We focus on things that matter, this way we are able to ensure easy maintenance and keep the system well-organized, well-planned and sweet.

* Elegance. `utron` uses golang best practises. We are not afraid of heights, it's just that we need a parachute in our backpack. The source code is heavily documented, any functionality should be well explained and well tested.

## Motivation
After two years of playing with golang, I have looked on some of my projects and asked myself: "How golang is that?"

So, `utron` is my reimagining of lightweight MVC, that maintains the golang spirit, and works seamlessly with the current libraries.


## Installation

	$ go get github.com/gernest/utron

## The MVC
There is nothing revolutionary about MVC that `utron` brings to the table.

* M is for models, they are the data structures that help with data persistence, utron uses [gorm](https://github.com/jinzhu/gorm) an existing Object Relational Mapper for golang. So if you are familiar with gorm then you are good on the M part.

* V is for Views. Views are templates that render the final output. `utron` uses golang standard templates. You don't have to learn anything new, just the text/template package to master views.

* C is for controllers. This is where the application logic lives. In order to achieve modularity, there are some things that utron requires of controllers. This subject is explained in more detail below.

With the power of composition and inheritance, `utron` achieves a beautiful MVC workflow. I recommend you read the source code, it is well documented so as to demystify any magical unicorns.


We will create a TODO List application in `utron` to explore all components that makes `utron` MVC tick. The source code of the final application is included in this repository and can be found here [utron todoMVC](fixtures/todo)

# TODO list application with `utron`

## Project structure
This is the structure of the `todo` list application that will showcase how you can build web apps with `utron`:

```shell
todo
├── config
│   ├── app.json
│   ├── app.toml
│   └── app.yml
├── controllers
│   └── todo.go
├── models
│   └── todo.go
├── static
│   └── todo.css
├── views
│   ├── error.html
│   └── index.html
└── main.go

5 directories, 9 files
```

I have included three configuration files to show how they work, but you are better off with just one.

## Configurations
utron support yaml, json and toml configurations files. In our todo app, we put the configuration files in the config directory. I have included all three formats for clarity, you can be just fine with either one of them.

`utron` searches for a file named `app.json`, or `app.yml` or `app.toml` in the config directory. The first to be found is the one to be used.

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
	"database_conn": "postgres://postgres:postgres@localhost/todo"
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
database      | DATABASE         | the name of the database you use, e.g. postgres, mysql, foundation
database_conn | DATABASE_CONN    | connection string to your database

If you haven't specified explicitly the location of the configuration directory, it defaults to the directory named `config` in the current working directory.

## Models
`utron` uses the [gorm](https://github.com/jinzhu/gorm) library as its Object Relational Mapper, so you won't need to learn anything fancy. In our todo app, we need to define a `Todo` model that will be used to store our todo details.

In the file `models/todo.go` we define our todo model like this

```go
package models

import (
	"time"

	"github.com/gernest/utron"
)

type Todo struct {
	ID        int       `schema: "-"`
	Body      string    `schema:"body"`
	CreatedAt time.Time `schema:"-"`
	UpdatedAt time.Time `schema:"-"`
}

func init() {
	utron.RegisterModels(&Todo{})
}
```

Notice that we need to register our model by calling `utron.RegisterModels(&Todo{})` in the `init` function otherwise `utron` won't be aware of the model.

`utron` will automatically create the table `todos` if it doesn't exist.

Don't be confused by the `schema` tag, I just added them since we will use the [schema](https://github.com/gorilla/schema) package to decode form values(this has nothing to do with `utron`, you can use whatever form library you fancy.)


## Controllers
`utron` controllers are structs that implement the `Controller` interface. To help make `utron` usable, `utron` provides a `BaseController` which implements the `Controller` interface and offers additional conveniences to help in composing reusable code.

You get all the benefits of `BaseController` by embedding it in your struct. Our `TODO` Controller is in the `controller/todo.go`

```go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gernest/utron"
	"github.com/gernest/utron/fixtures/todo/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type TODO struct {
	*utron.BaseController
	Routes []string
}

func (t *TODO) Home() {
	todos := []*models.Todo{}
	t.Ctx.DB.Order("created_at desc").Find(&todos)
	t.Ctx.Data["List"] = todos
	t.Ctx.Template = "index"
	t.HTML(http.StatusOK)
}
func (t *TODO) Create() {
	todo := &models.Todo{}
	req := t.Ctx.Request()
	req.ParseForm()
	if err := decoder.Decode(todo, req.PostForm); err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}

	t.Ctx.DB.Create(todo)
	t.Ctx.Redirect("/", http.StatusFound)
}

func (t *TODO) Delete() {
	todoID := t.Ctx.Params["id"]
	ID, err := strconv.Atoi(todoID)
	if err != nil {
		t.Ctx.Data["Message"] = err.Error()
		t.Ctx.Template = "error"
		t.HTML(http.StatusInternalServerError)
		return
	}
	t.Ctx.DB.Delete(&models.Todo{ID: ID})
	t.Ctx.Redirect("/", http.StatusFound)
}

func NewTODO() *TODO {
	return &TODO{
		Routes: []string{
			"get;/;Home",
			"post;/create;Create",
			"get;/delete/{id};Delete",
		},
	}
}

func init() {
	utron.RegisterController(NewTODO())
}
```

Note that we registered our controller by calling `utron.RegisterController(NewTODO())` in the `init` function
so as to make `utron` aware of our controller. See Routing section below for more explanation of what the controller is doing.


## Routing

By registering a controller, there are two ways of assigning routes.

### case 1- vanilla routing
By registering a Controller, routes are auto-generated for the controller methods. The format is `/:controller/:method` where `:controller` is the lowercase name of the Controller, and `:method` is its method in lowercase.

so `(*TODO) Hello()` will map to `/todo/hello`

### case 2- Specifying Routes field
The user controller can define a field named `Routes` it should be of type `[]string`, then you can assign routes by appending route string to the Routes field.

This is a better explanation from comments on the `router.go` file.

```go
		// if there is any field named Routes, and it is of signature []string
		// then the field's value is used to overide the patterns defined earlier.
		//
		// It is not necessary for every user implementation to define method named Routes
		// If we cant find it then we just ignore its use( fallback to defaults).
		//
		// Route strings, are of the form "httpMethods;path;method"
		// where httpMethods: is a comma separated list of http method strings
		//                  e.g GET,POST,PUT.
		//                  The case does not matter, you can use lower case or upper case characters
		//                  or even mixed case, that is get,GET,gET and GeT will all be treated as GET
		//
		//        path:     Is a url path or pattern, utron uses the gorilla mux package. So, everything you can do
		//                  with gorilla mux url path then you can do here.
		//                  e.g /hello/{world}
		//                  Don't worry about the params, they will be accessible via .Ctx.Params field in your
		//                  controller.
		//
		//        method:   The name of the user Controller method to execute for this route.
```

So, that explains the following lines in our `todo` app in `controllers/todo.go`

```go
func NewTODO() *TODO {
	return &TODO{
		Routes: []string{
			"get;/;Home",
			"post;/create;Create",
			"get;/delete/{id};Delete",
		},
	}
}
```
### case 3: using routes file
You can define routes in a file in the `config` directory. The supported formats are json, toml and yaml.

`utron` will look for file named `routes.json`, `routes.toml` or `routes.yml` in that order, the first to be found is the one to be used.

I have included a sample routes file in [fixtures/config/routes.json](fixtures/config/routes.json).

The difference with case 2 above is you will need to specify the name of the controller explicitly. That is for `TODO` controller, we can define the home route string in routes file like `get;/;TODO.Home`.

We won't use this in our TODO list app, but you can find it useful in your use case.'

## Views
`utron` views are golang templates. This is the content of `views/index.html`:

```html
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>Utron Todo MVC</title>
	<link href="/static/todo.css" rel="stylesheet">
</head>
<body>
<form method="post" action="/create">
    <table>
        <tr>
           <th>
               Create A TODO
           </th>
        </tr>
        <tr>
            <td>
                <input name="body">
            </td>
            <td>
                <button type="submit">create</button>
            </td>
        </tr>
    </table>
</form>
<table>
    <tr>
        <th>
            My TODO LIST
        </th>
    </tr>
    {{range $k,$v:=.List}}
    <tr>
        <td>
            {{$v.ID}}
        </td>
        <td>
            {{$v.Body}}
        </td>
        <td>
            <a href="/delete/{{$v.ID}}">
                <button>Delete</button>
            </a>
        </td>
    </tr>
    {{end}}
</table>
</body>
</html>
```


Note that we have access to `.List` in our view. This is set in the controller, additionally you can access the application configuration via `.Config` context.

Above is a simple golang template to render our `todo` list application.

## The main.go file

```go
package main

import (
	"github.com/gernest/utron"
	_ "github.com/gernest/utron/fixtures/todo/controllers"
	_ "github.com/gernest/utron/fixtures/todo/models"
)

func main() {
	utron.Run()
}
```

## Running the TODO app
In case you want to run the app we just created, it is included in this repository in [fixtures/todo](fixtures/todo)

* Prerequisite
 - a working database connection (postres, mysql or foundation)
 - golang toolchain installed and the `go` command in your system $PATH.

step 1 Install `utron` which will also include the todo app

	$ go get github.com/gernest/utron

step 2 cd into the todo app directory

	$ cd $GOPATH/src/github.com/gernest/utron/fixtures/todo

step 3 install dependency

	$ go get github.com/gorilla/schema

step 4 edit `config/app.json` by setting database and database_conn to your values

step 5 run the app

	go run main.go

If you see something like this

	$ 2015/09/15 18:27:24 >>INFO>> starting server at http://localhost:8090

Then everything is okay, open `http://localhost:8090` in your browser to start writing your todos.
If you experience anything different, redo the steps and make sure you did them in order and with no errors. If so, and it still doesn't work, open an [issue](https://github.com/gernest/utron/issues).

## Screenshot
![todo app with utron](fixtures/todo.png)

# Contributing

Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest <geofreyernest@live.com>

Twitter  : [@gernesti](https://twitter.com/gernesti)

Facebook : [Geofrey Ernest](https://www.facebook.com/geofrey.ernest.35)

# Are you hiring?

I have 2 years experience working with golang and 5 years of web development experience. I don't have a juicy CV because I'm in Tanzania, and almost everyone I graduated with has no job, has never had one and doubts he/she will ever land one (unless, of course, you are lucky which I'm not.)

If you are a recruiter, please get in touch, because I have my full trust that code will get me out of this abysmal poverty.

# Acknowledgements
These amazing projects have made `utron` possible:

* [gorilla mux](https://github.com/gorilla/mux)
* [ita](https://github.com/gernest/ita)
* [gorm](https://github.com/jinzhu/gorm)
* [alice](https://github.com/justinas/alice)
* [golang](http://golang.org)


# Roadmap

*  Fix a lot of typos (English is my third language).

# Licence

This project is released under the MIT licence. See [LICENCE](LICENCE) for more details.
