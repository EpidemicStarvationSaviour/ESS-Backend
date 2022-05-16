package user

type Role int

const (
	NoLogin   Role = 0
	Supplier  Role = 1
	Rider     Role = 2
	Purchaser Role = 3
	Leader    Role = 4
	Admin     Role = 5
	SysAdmin  Role = 6
)
