package types

type AccountType = string

const (
	AdminAccount   AccountType = "admin"
	SellerAccount  AccountType = "seller"
	RegularAccount AccountType = "regular"
	AnyAccount     AccountType = ""
)
