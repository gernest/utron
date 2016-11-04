# models
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

