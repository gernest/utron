# models
--
    import "github.com/NlaakStudios/gowaf/models"


## Usage

```go
const (

	// UserStateVerifyEmailSend -> Need to send out email verification
	UserStateVerifyEmailSend = 1
	// UserStateVerifyEmailSent -> Email verification has been sent
	UserStateVerifyEmailSent = 2
	// UserStateVerifyEmailDone -> User clicked on link in email, Verified, Account Active
	UserStateVerifyEmailDone = 3
	// UserStateBanned -> Account has been banned. No access granted other than guest
	UserStateBanned = 4
	// UserStateIdle -> User Signed in but no activity for 5 min (300 sec)
	UserStateIdle = 5
	// UserStateSignedIn -> User is Signed In
	UserStateSignedIn = 6
	// UserStateSignedOut -> User is Signed Out
	UserStateSignedOut = 7

	// UserAccessGuest -> default, no account, Guest access
	UserAccessGuest = 0
	// UserAccessMember -> Active account with access to Member content
	UserAccessMember = 1
	// UserAccessEmployee -> Active account with access to Member and Employee content
	UserAccessEmployee = 2
	// UserAccessAdmin -> Active account with access to Member, Employee and Admin content
	UserAccessAdmin = 3
)
```

```go
const (
	//PhoneTypeUnknown represents a defaul unknown phone type
	PhoneTypeUnknown = byte(0)
	//PhoneTypeMobile represents a Mobile or Cell phone number
	PhoneTypeMobile = byte(1)
	//PhoneTypeHome represents a home phone number
	PhoneTypeHome = byte(2)
	//PhoneTypeBusiness represents a business phone number
	PhoneTypeBusiness = byte(3)
	//PhoneTypeFax represents a Fax phone number
	PhoneTypeFax = byte(4)
)
```

```go
const (
	//KYCNone used in UserKYCStruct.Tier to represent no KYC verification exists.
	KYCNone = 0
	//KYCBasic used in UserKYCStruct.Tier to represent basic verification only.
	KYCBasic = 1
	//KYCIntermediate used in UserKYCStruct.Tier to represent intermediate level of verification.
	KYCIntermediate = 2
	//KYCAdvanced used in UserKYCStruct.Tier to represent Advanced for Full level of verification.
	KYCAdvanced = 3
)
```

#### func  CoinMap

```go
func CoinMap() map[string]interface{}
```
CoinMap Converts CoinStruct to a map[string]interface{}

#### func  LoadCoinInfo

```go
func LoadCoinInfo()
```
LoadCoinInfo loads coin info from config.json file

#### func  SaveCoinInfo

```go
func SaveCoinInfo()
```
SaveCoinInfo saves current coin info to config.json file

#### func  SetCoinInfo

```go
func SetCoinInfo(coin CoinStruct)
```
SetCoinInfo set the default (from consts) for a coin

#### type Account

```go
type Account struct {
	ID             int       `schema:"id"`
	CreatedAt      time.Time `schema:"created"`
	UpdatedAt      time.Time `schema:"updated"`
	Username       string    `valid:"required,length(6|16)" schema:"username"`
	Email          string    `valid:"required,length(6|16)" schema:"email"`
	Password       string    `gorm:"-" valid:"required,length(6|24)" schema:"password"`
	VerifyPassword string    `gorm:"-" valid:"required,length(6|24)" schema:"verify_password"`
	HashedPassword string    `schema:"hashed_password"`
	State          byte      `schema:"state"`
	Access         byte      `schema:"access"`
}
```

Account is used to represent a user for authentication

#### func (*Account) CheckPassword

```go
func (m *Account) CheckPassword(dbHash, givenPW string) bool
```
CheckPassword checks that the password hash in the database matches the password
the user just gave. Return TRUE if valid

#### func (*Account) HTMLForm

```go
func (m *Account) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Account) HTMLView

```go
func (m *Account) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Account) MultiLine

```go
func (m *Account) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model {Username}:
{Person.SingleLine()} {Email.Address} {Company.SingleLine()}

#### func (*Account) SetPassword

```go
func (m *Account) SetPassword(pw string) string
```
SetPassword create a password hash

#### func (*Account) SingleLine

```go
func (m *Account) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model
{Username}: {Email.Address} [{ID},{CID},{PID}]}

#### func (*Account) Validate

```go
func (m *Account) Validate() error
```
Validate is used to verifiy password hash match

#### type Address

```go
type Address struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Address1  string    `valid:"required" schema:"address1"`
	Address2  string    `valid:"required" schema:"address2"`
	City      string    `valid:"required" schema:"city"`
	State     string    `valid:"required" schema:"state"`
	Zip       string    `valid:"required" schema:"zip"`
	County    string    `valid:"required" schema:"county"`
	Country   string    `valid:"required" schema:"country"`
}
```

Address contains general Address

#### func (*Address) HTMLForm

```go
func (m *Address) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Address) HTMLView

```go
func (m *Address) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Address) IsValid

```go
func (m *Address) IsValid() error
```
IsValid returns error if address is not complete

#### func (*Address) MultiLine

```go
func (m *Address) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Address) SingleLine

```go
func (m *Address) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type CoinStruct

```go
type CoinStruct struct {
	Name      string
	Symbol    string
	Decimals  uint
	Supply    CoinSupplyStruct
	Company   Company
	Compliant CompliantStruct
	URLS      CoinURLStruct
}
```

CoinStruct stores all Coin information including Company, Compliantcy

```go
var CoinSettings CoinStruct
```
CoinSettings Holds all information pertaining to the Coin

#### func  GetCoinInfo

```go
func GetCoinInfo() CoinStruct
```
GetCoinInfo Gets the Coin to User defined Data

#### type CoinSupplyStruct

```go
type CoinSupplyStruct struct {
	ICO   uint64
	Dev   uint64
	OAM   uint64
	PLT   uint64
	Total uint64
}
```

CoinSupplyStruct stores the Total Coin Supply aand the breakdown

#### type CoinURLStruct

```go
type CoinURLStruct struct {
	LandingPage   string
	DashboardPage string
	API           string
	WhitePaper    string
}
```

CoinURLStruct stores pre-built URLS for the Coin ICO

#### type Company

```go
type Company struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Name      string    `schema:"name"`
	ContactID int       `schema:"contact_id"`
	Person    Person    `gorm:"foreignkey:ContactID"`
	PhoneID   int       `schema:"phone_id"`
	Phone     Phone     `gorm:"foreignkey:PhoneID"`
	FaxID     int       `schema:"fax_id"`
	Fax       Phone     `gorm:"foreignkey:FaxID"`
}
```

Company stores information about the company

#### func (*Company) HTMLForm

```go
func (m *Company) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Company) HTMLView

```go
func (m *Company) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Company) MultiLine

```go
func (m *Company) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Company) SingleLine

```go
func (m *Company) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type CompliantStruct

```go
type CompliantStruct struct {
	AML string
	KYC string
}
```

CompliantStruct stores information about AML and KYC states

#### type Email

```go
type Email struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Address   string    `schema:"address"`  //bob1234@gmail.com
	Username  string    `schema:"username"` //bob1234
	Domain    string    `schema:"domain"`   //gmail.com
}
```

Email contains a breakdown of a email TODO: Update to use/integrate "net/mail"
and Address

#### func (*Email) HTMLForm

```go
func (m *Email) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Email) HTMLView

```go
func (m *Email) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Email) MultiLine

```go
func (m *Email) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Email) SingleLine

```go
func (m *Email) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type Example

```go
type Example struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Number    int       `valid:"required" schema:"number"` //123456
	String    string    `valid:"required" schema:"string"` //Test String
	Toggle    bool      `valid:"required" schema:"toggle"` //True/False, On/Off, Yes/No
	Float     float64   `valid:"required" schema:"float"`  //$12.56, 1,123,456.987654321
}
```

Example contains general Example

#### func (*Example) HTMLForm

```go
func (m *Example) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Example) HTMLView

```go
func (m *Example) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Example) MultiLine

```go
func (m *Example) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Example) SingleLine

```go
func (m *Example) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type Gender

```go
type Gender struct {
	ID         int       `schema:"id"`
	CreatedAt  time.Time `schema:"created"`
	UpdatedAt  time.Time `schema:"updated"`
	ClaimedSex string    `schema:"claimed_sex"` // what they claim -> male, female, gay, lesbian, transgender, etc
	BioSex     byte      `schema:"legal_sex"`   //What is on birth certificate / under the hood? 0=Unknown, 1=Male, 2=Female
}
```

Gender aims to be LGBT+ compliant and is primarly used for referencing the
'Person' in the webapp and templating system

#### func (*Gender) BioSexToString

```go
func (m *Gender) BioSexToString(gender byte) string
```
BioSexToString translates the byte value to human readable friendly string

#### func (*Gender) HTMLForm

```go
func (m *Gender) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Gender) HTMLView

```go
func (m *Gender) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Gender) IsValid

```go
func (m *Gender) IsValid() error
```

#### func (*Gender) MultiLine

```go
func (m *Gender) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Gender) SingleLine

```go
func (m *Gender) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type Model

```go
type Model struct {
	*gorm.DB
}
```

Model facilitate database interactions, supports postgres, mysql and foundation

#### func  NewModel

```go
func NewModel() *Model
```
NewModel returns a new Model without opening database connection

#### func (*Model) AutoMigrateAll

```go
func (m *Model) AutoMigrateAll()
```
AutoMigrateAll runs migrations for all the registered models

#### func (*Model) Count

```go
func (m *Model) Count() int
```
Count returns the number of registered models

#### func (*Model) IsOpen

```go
func (m *Model) IsOpen() bool
```
IsOpen returns true if the Model has already established connection to the
database

#### func (*Model) OpenWithConfig

```go
func (m *Model) OpenWithConfig(cfg *config.Config) error
```
OpenWithConfig opens database connection with the settings found in cfg

#### func (*Model) Register

```go
func (m *Model) Register(values ...interface{}) error
```
Register adds the values to the models registry

#### type ModelStats

```go
type ModelStats struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	//Total is the total number of models in existence. Avaiable in temaplate as {{.models_total}}
	Total uint `schema:"total"`
	//References is how many times per day the model is referenced.Avaiable in temaplate as {{.model_references}}
	References uint `schema:"referenced"`
	//Is the rounded, total percent of models that are active. Avaiable in temaplate as {{.model_pct_active}}
	PctActive uint `schema:"pct_active"`
	//Unused is the total number of models that are not being referenced. Avaiable in temaplate as {{.model_unused}}
	Unused uint `schema:"unused"`
	//Active is the total number of models that are being referenced. Avaiable in temaplate as {{.model_active}}
	Active uint `schema:"active"`
	//Archived is the total number of models that are currently archived. Avaiable in temaplate as {{.model_archived}}
	Archived uint `schema:"archived"`
}
```

ModelStats holds various stats for each model int he database and can be
displayed in the models dashboard

#### func (*ModelStats) NewModelStats

```go
func (m *ModelStats) NewModelStats()
```

#### type Note

```go
type Note struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	PersonID  int       `schema:"person_id"`
	Person    Person    `gorm:"foreignkey:PersonID"`
	Body      string    `schema:"body"`
}
```

Note

#### func (*Note) HTMLForm

```go
func (m *Note) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Note) HTMLView

```go
func (m *Note) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Note) MultiLine

```go
func (m *Note) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Note) SingleLine

```go
func (m *Note) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type Person

```go
type Person struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Dob       time.Time `schema:"dob"`

	//GenderID is the UID of the Person's Gender as found in the gender table
	GenderID int    `schema:"gender_id"`
	Gender   Gender `gorm:"foreignkey:GenderID"`

	//NameID is the UID of the Person's name as found in the person_name table
	NameID     int        `schema:"name_id"`
	PersonName PersonName `gorm:"foreignkey:NameID"`

	//EmailID is the UID of the Person's email as found in the email table
	EmailID int   `schema:"email_id"`
	Email   Email `gorm:"foreignkey:EmailID"`

	//TypeID is the UID of the Person's Type as found in the person_type table
	TypeID     int        `schema:"type_id"`
	PersonType PersonType `gorm:"foreignkey:TypeID"`

	//PhoneID is the UID of the Person's Phone info as found in the phone table
	PhoneID int   `schema:"phone_id"`
	Phone   Phone `gorm:"foreignkey:PhoneID"`
}
```

Person contains all data pertaining to a individual person

#### func (*Person) HTMLForm

```go
func (m *Person) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of a Person Model

#### func (*Person) HTMLView

```go
func (m *Person) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of a Person Model

#### func (*Person) MultiLine

```go
func (m *Person) MultiLine() string
```
MultiLine returns a formatted multi-line text representing a Person Model

#### func (*Person) SingleLine

```go
func (m *Person) SingleLine() string
```
SingleLine returns a formatted single line text representing a Person Model

#### type PersonName

```go
type PersonName struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Prefix    string    `schema:"prefix"`  //ie. Mr
	First     string    `schema:"first"`   //William
	Middle    string    `schema:"middle"`  //Blaine
	Last      string    `schema:"last"`    //Doe
	Suffix    string    `schema:"suffix"`  //Sr
	GoesBy    string    `schema:"goes_by"` //Bob
}
```

PersonName hold the complete name of a person

#### func (*PersonName) HTMLForm

```go
func (m *PersonName) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*PersonName) HTMLView

```go
func (m *PersonName) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*PersonName) MultiLine

```go
func (m *PersonName) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*PersonName) SingleLine

```go
func (m *PersonName) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type PersonType

```go
type PersonType struct {
	ID        int       `schema:"id"`
	CreatedAt time.Time `schema:"created"`
	UpdatedAt time.Time `schema:"updated"`
	Name      string    `schema:"name"`
}
```

PersonType provides a list of avilable "title" such as: 0 Unknown 1 Adjuster 2
Property Owner 3 Attorney 4 Paralegal 5 Contractor

#### func (*PersonType) HTMLForm

```go
func (m *PersonType) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*PersonType) HTMLView

```go
func (m *PersonType) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*PersonType) MultiLine

```go
func (m *PersonType) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*PersonType) SingleLine

```go
func (m *PersonType) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type Phone

```go
type Phone struct {
	ID          int       `schema:"id"`
	CreatedAt   time.Time `schema:"created"`
	UpdatedAt   time.Time `schema:"updated"`
	CountryCode string    `schema:"code"`
	AreaCode    string    `schema:"area"`
	Number      string    `schema:"number"`
	PhoneType   byte      `schema:"phone_type"`
}
```

PhoneStruct is used to breakdown and store phone numbers

#### func (*Phone) HTMLForm

```go
func (m *Phone) HTMLForm() string
```
HTMLForm returns a HTML5 code representing a form of the Model

#### func (*Phone) HTMLView

```go
func (m *Phone) HTMLView() string
```
HTMLView returns a HTML5 code representing a view of the Model

#### func (*Phone) MultiLine

```go
func (m *Phone) MultiLine() string
```
MultiLine returns a formatted multi-line text representing the Model

#### func (*Phone) PhoneTypeToString

```go
func (m *Phone) PhoneTypeToString(pt byte) string
```
PhoneTypeToString given a valid PhoneType Byte value will return the string
representation

#### func (*Phone) SingleLine

```go
func (m *Phone) SingleLine() string
```
SingleLine returns a formatted single line text representing the Model

#### type UserKYCSignature

```go
type UserKYCSignature struct {
	Verified bool
	//UID of Person Node/Record who Verified
	Person uint64
	//The Time in which verification occured
	When uint64
}
```

UserKYCSignature Hold Verification status as well as who verified and when

#### type UserKYCStruct

```go
type UserKYCStruct struct {
	//User Current verified Tier Level (0..3)
	//
	Tier  uint
	Tier1 UserKYCTierOneStruct
	Tier2 UserKYCTierTwoStruct
	Tier3 UserKYCTierThreeStruct
}
```

UserKYCStruct contains users kyc "Know Your Customer" extended data

#### type UserKYCTierOneStruct

```go
type UserKYCTierOneStruct struct {
	//Person is the UID of a Person Node/Record
	Person uint64
	//Address is the UID of a Address Node/Record
	Address uint64
	//Phone is the UID of a Phone Node/Record
	Phone uint64
	//NatID is a countries National Identification Number (SSN for USA) for the person
	NatID string
	//PictureOfIDPath is the Full Path and filename to the PNG containing a cropped scan of Picture Identification
	PictureOfIDPath string
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
}
```

UserKYCTierOneStruct Basic Customer (User) information (Tier 1)

#### type UserKYCTierThreeStruct

```go
type UserKYCTierThreeStruct struct {
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
}
```

UserKYCTierThreeStruct Advanced Customer (User) information (Tier 3)

#### type UserKYCTierTwoStruct

```go
type UserKYCTierTwoStruct struct {
	//Signature contains the information about who and when verification took place.
	Signature UserKYCSignature
}
```

UserKYCTierTwoStruct Intermediate Customer (User) information (Tier 2)

#### type UserSecurityStruct

```go
type UserSecurityStruct struct {
	//TwoFA enable / diable Two Factor Authentication
	TwoFA bool
	//EmailOnLogin enable sending email for each account login
	EmailOnLogin bool
	//EmailWithdrawConfirm enable sending email to confirm withdrawl
	EmailWithdrawConfirm bool
	//UseOfflineWallet enable using offline wallet for cold storage
	UseOfflineWallet bool
	//AutoTransferApproved transfer required funds from offline wallet on approved/verified transaction
	AutoTransferApproved bool
	//MaxCoinsOnlineWallet maximum number of coins to keep in online wallet.
	//Everything over will be transferred to offline wallet
	MaxCoinsOnlineWallet uint32
}
```

UserSecurityStruct Contains users Security preferences

#### type Wallet

```go
type Wallet struct {
	UserID uint64 //AccountStruct|dGraph Node UID? Who Owns this wallet...For authentication

	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}
```

Wallet stores private and public keys
