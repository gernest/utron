# controller
--
    import "github.com/NlaakStudios/gowaf/controller"


## Usage

```go
var Decoder = schema.NewDecoder()
```
Decoder is use to decode the schema

#### func  GetCtrlFunc

```go
func GetCtrlFunc(ctrl Controller) func() Controller
```
GetCtrlFunc returns a new copy of the contoller everytime the function is called

#### type Account

```go
type Account struct {
	BaseController
	Routes []string
}
```

Account is the controller for the Account Model

#### func (*Account) Index

```go
func (a *Account) Index()
```
Index displays the account landing (index) page

#### func (*Account) Login

```go
func (a *Account) Login()
```
Login displays login page for GET and processes on POST

#### func (*Account) Logout

```go
func (a *Account) Logout()
```
Logout logs the user out of they are logged in

#### func (*Account) Register

```go
func (a *Account) Register()
```
Register displays registration page for GET and processes form data on POST

#### func (*Account) SendEmailVerification

```go
func (a *Account) SendEmailVerification(acct models.Account)
```
SendEmailVerification Sends a Verification Email to the user registering

#### type Address

```go
type Address struct {
	BaseController
	Routes []string
}
```

Account is the controller for the Account Model

#### func (*Address) Create

```go
func (c *Address) Create()
```
Create creates a Address item

#### func (*Address) Delete

```go
func (c *Address) Delete()
```
Delete deletes a Address item

#### func (*Address) Edit

```go
func (c *Address) Edit()
```

#### func (*Address) Index

```go
func (c *Address) Index()
```
Home renders a Address list

#### func (*Address) View

```go
func (c *Address) View()
```
Delete deletes a Address item

#### type BaseController

```go
type BaseController struct {
	//Stats  *models.ModelStats
	Ctx    *base.Context
	Routes []string
}
```

BaseController implements the Controller interface, It is recommended all user
defined Controllers should embed *BaseController.

#### func (*BaseController) HTML

```go
func (b *BaseController) HTML(code int)
```
HTML renders text/html with the given code as status code

#### func (*BaseController) JSON

```go
func (b *BaseController) JSON(code int)
```
JSON renders application/json with the given code

#### func (*BaseController) New

```go
func (b *BaseController) New(ctx *base.Context)
```
New sets ctx as the active context

#### func (*BaseController) Render

```go
func (b *BaseController) Render() error
```
Render commits the changes made in the active context.

#### func (*BaseController) RenderJSON

```go
func (b *BaseController) RenderJSON(value interface{}, code int)
```
RenderJSON encodes value into json and renders the response as JSON

#### func (*BaseController) String

```go
func (b *BaseController) String(code int)
```
String renders text/plain with given code as status code

#### type Company

```go
type Company struct {
	BaseController
	Routes []string
}
```

Company is a controller for Company list

#### func (*Company) Create

```go
func (c *Company) Create()
```
Create creates a Company item

#### func (*Company) Delete

```go
func (c *Company) Delete()
```
Delete deletes a Company item

#### func (*Company) Index

```go
func (c *Company) Index()
```
Home renders a Company list

#### func (*Company) View

```go
func (c *Company) View()
```
Delete deletes a Company item

#### type Controller

```go
type Controller interface {
	New(*base.Context)
	Render() error
}
```

Controller is an interface for gowaf controllers

#### func  NewAccount

```go
func NewAccount() Controller
```
NewAccount returns a new account controller object

#### func  NewAddress

```go
func NewAddress() Controller
```
NewAddress returns a new Address list controller

#### func  NewCompany

```go
func NewCompany() Controller
```
NewCompany returns a new Company list controller

#### func  NewEmail

```go
func NewEmail() Controller
```
NewEmail returns a new Email list controller

#### func  NewExample

```go
func NewExample() Controller
```
NewExample returns a new account controller object

#### func  NewGender

```go
func NewGender() Controller
```
NewGender returns a new Gender list controller

#### func  NewLanding

```go
func NewLanding() Controller
```
NewLanding returns a new account controller object

#### func  NewNote

```go
func NewNote() Controller
```
NewNote returns a new Note list controller

#### func  NewPerson

```go
func NewPerson() Controller
```
NewPerson returns a new Person list controller

#### func  NewPersonName

```go
func NewPersonName() Controller
```
NewPersonName returns a new PersonName list controller

#### func  NewPersonType

```go
func NewPersonType() Controller
```
NewPersonType returns a new PersonType list controller

#### func  NewPhone

```go
func NewPhone() Controller
```
NewPhone returns a new Phone list controller

#### type Email

```go
type Email struct {
	BaseController
	Routes []string
}
```

Email is a controller for Email list

#### func (*Email) Create

```go
func (c *Email) Create()
```
Create creates a Email item

#### func (*Email) Delete

```go
func (c *Email) Delete()
```
Delete deletes a Email item

#### func (*Email) Edit

```go
func (c *Email) Edit()
```

#### func (*Email) Index

```go
func (c *Email) Index()
```
Home renders a Email list

#### func (*Email) View

```go
func (c *Email) View()
```
Delete deletes a Email item

#### type Example

```go
type Example struct {
	BaseController
	Routes []string
}
```

Example is the controller for the Example Model

#### func (*Example) Create

```go
func (a *Example) Create()
```
Create creates a new model in the database

#### func (*Example) Delete

```go
func (a *Example) Delete()
```
Delete deletes a model in the database with correct access level

#### func (*Example) Index

```go
func (a *Example) Index()
```
Index displays the account example (index) page

#### func (*Example) List

```go
func (a *Example) List()
```
List shows a paginated list of all model items based on filter / search info

#### func (*Example) ViewEdit

```go
func (a *Example) ViewEdit()
```
Edit edits an existing model in the database with correct access level

#### type Gender

```go
type Gender struct {
	BaseController
	Routes []string
}
```

Gender is a controller for Gender list

#### func (*Gender) Create

```go
func (c *Gender) Create()
```
Create creates a Gender item

#### func (*Gender) Delete

```go
func (c *Gender) Delete()
```
Delete deletes a Gender item

#### func (*Gender) Edit

```go
func (c *Gender) Edit()
```

#### func (*Gender) Index

```go
func (c *Gender) Index()
```
Home renders a Gender list

#### func (*Gender) View

```go
func (c *Gender) View()
```
Delete deletes a Gender item

#### type Landing

```go
type Landing struct {
	BaseController
	Routes []string
}
```

Landing is the controller for the Landing Model

#### func (*Landing) Contact

```go
func (a *Landing) Contact()
```
Contact implements as Contact us form

#### func (*Landing) Index

```go
func (a *Landing) Index()
```
Index displays the account landing (index) page

#### func (*Landing) Services

```go
func (a *Landing) Services()
```
Services displays the services page

#### type Note

```go
type Note struct {
	BaseController
	Routes []string
}
```

Note is a controller for Note list

#### func (*Note) Create

```go
func (t *Note) Create()
```
Create creates a Note item

#### func (*Note) Delete

```go
func (t *Note) Delete()
```
Delete deletes a Note item

#### func (*Note) Index

```go
func (t *Note) Index()
```
Home renders a Note list

#### type Person

```go
type Person struct {
	BaseController
	Routes []string
}
```

Person is a controller for Person list

#### func (*Person) Create

```go
func (c *Person) Create()
```
Create creates a Person item

#### func (*Person) Delete

```go
func (c *Person) Delete()
```
Delete deletes a Person item

#### func (*Person) Index

```go
func (c *Person) Index()
```
Home renders a Person list

#### func (*Person) View

```go
func (c *Person) View()
```
Delete deletes a Person item

#### type PersonName

```go
type PersonName struct {
	BaseController
	Routes []string
}
```

PersonName is a controller for PersonName list

#### func (*PersonName) Create

```go
func (c *PersonName) Create()
```
Create creates a PersonName item

#### func (*PersonName) Delete

```go
func (c *PersonName) Delete()
```
Delete deletes a PersonName item

#### func (*PersonName) Index

```go
func (c *PersonName) Index()
```
Home renders a PersonName list

#### func (*PersonName) View

```go
func (c *PersonName) View()
```
Delete deletes a PersonName item

#### type PersonType

```go
type PersonType struct {
	BaseController
	Routes []string
}
```

PersonType is a controller for PersonType list

#### func (*PersonType) Create

```go
func (c *PersonType) Create()
```
Create creates a PersonType item

#### func (*PersonType) Delete

```go
func (c *PersonType) Delete()
```
Delete deletes a PersonType item

#### func (*PersonType) Index

```go
func (c *PersonType) Index()
```
Home renders a PersonType list

#### func (*PersonType) View

```go
func (c *PersonType) View()
```
Delete deletes a PersonType item

#### type Phone

```go
type Phone struct {
	BaseController
	Routes []string
}
```

Phone is a controller for Phone list

#### func (*Phone) Create

```go
func (c *Phone) Create()
```
Create creates a Phone item

#### func (*Phone) Delete

```go
func (c *Phone) Delete()
```
Delete deletes a Phone item

#### func (*Phone) Index

```go
func (c *Phone) Index()
```
Home renders a Phone list

#### func (*Phone) View

```go
func (c *Phone) View()
```
Delete deletes a Phone item
