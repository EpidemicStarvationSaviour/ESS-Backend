package response

// ResponseType
type Type int

const TypeJSON Type = 1
const TypeFile Type = 2
const TypeFailed Type = 3 // http status code
const TypeRedirect Type = 4
const TypeImage Type = 5

// for Debug
func (t Type) String() string {
	switch t {
	case TypeJSON:
		return "JSON"
	case TypeFile:
		return "File"
	case TypeFailed:
		return "Failed"
	case TypeRedirect:
		return "Redirect"
	case TypeImage:
		return "Images"
	default:
		return "<unknown>"
	}
}
