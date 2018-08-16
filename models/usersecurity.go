package models

// UserSecurityStruct Contains users Security preferences
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
