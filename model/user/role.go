package user

type Role int

const (
	NoLogin  Role = 0
	SysAdmin Role = 1
	Admin    Role = 2
	EndUser  Role = 3
)
