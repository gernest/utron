package models

import (
	"github.com/fatih/structs"
)

//CoinSupplyStruct stores the Total Coin Supply aand the breakdown
type CoinSupplyStruct struct {
	ICO   uint64
	Dev   uint64
	OAM   uint64
	PLT   uint64
	Total uint64
}

//CompliantStruct stores information about AML and KYC states
type CompliantStruct struct {
	AML string
	KYC string
}

//CoinURLStruct stores pre-built URLS for the Coin ICO
type CoinURLStruct struct {
	LandingPage   string
	DashboardPage string
	API           string
	WhitePaper    string
}

//CoinStruct stores all Coin information including Company, Compliantcy
type CoinStruct struct {
	Name      string
	Symbol    string
	Decimals  uint
	Supply    CoinSupplyStruct
	Company   Company
	Compliant CompliantStruct
	URLS      CoinURLStruct
}

//CoinSettings Holds all information pertaining to the Coin
var CoinSettings CoinStruct

//SetCoinInfo set the default (from consts) for a coin
func SetCoinInfo(coin CoinStruct) {
	CoinSettings = coin
}

//GetCoinInfo Gets the Coin to User defined Data
func GetCoinInfo() CoinStruct {
	return CoinSettings
}

//SaveCoinInfo saves current coin info to config.json file
func SaveCoinInfo() {}

//LoadCoinInfo loads coin info from config.json file
func LoadCoinInfo() {}

//CoinMap Converts CoinStruct to a map[string]interface{}
func CoinMap() map[string]interface{} {
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	return structs.Map(CoinSettings)
}
