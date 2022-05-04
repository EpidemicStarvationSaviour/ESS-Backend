package response

// http status 30x redirect
func Redirect(code int, url string) Response {
	return Response{
		Type:         TypeRedirect,
		RedirectCode: code,
		RedirectURL:  url,
	}
}
