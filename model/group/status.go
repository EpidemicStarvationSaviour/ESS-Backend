package group

type Status int

const (
	All        Status = 0
	Created    Status = 1
	Submitted  Status = 2
	Delivering Status = 3
	Finished   Status = 4
)
