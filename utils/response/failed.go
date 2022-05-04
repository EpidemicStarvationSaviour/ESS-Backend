package response

// return http error status code
func Failed(code int) Response {
	return Response{
		Type:       TypeFailed,
		FailedCode: code,
	}
}
