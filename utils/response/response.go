package response

import (
	"encoding/json"
	"ess/utils/jsonu"
	"ess/utils/logging"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Type         Type
	Json         JSONResponse
	File         []byte // file data
	FileName     string
	FailedCode   int // http status code
	RedirectURL  string
	RedirectCode int
}

func (r *Response) Write(c *gin.Context) {
	switch r.Type {
	case TypeJSON:
		marshal := jsonu.Marshal(r.Json)
		c.JSON(http.StatusOK, json.RawMessage(marshal))
	case TypeFile:
		escape := url.QueryEscape(r.FileName)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", escape))
		c.Data(http.StatusOK, "application/octet-stream", r.File)
	case TypeRedirect:
		c.Redirect(r.RedirectCode, r.RedirectURL)
	case TypeFailed:
		c.Status(r.FailedCode)

	case TypeImage:
		c.File(r.FileName)
	}
}

func (r *Response) Logging() error {
	logging.Info(r.String())
	return nil
}

// summarize response into string for logging
func (r *Response) String() string {
	switch r.Type {
	case TypeJSON:
		return fmt.Sprintf("%+v", struct {
			Type Type
			Json string
		}{r.Type, r.ToJSON()})
	case TypeFile:
		return fmt.Sprintf("%+v", struct {
			Type     Type
			File     string // with only file size
			FileName string
		}{r.Type, fmt.Sprintf("<%d byte>", len(r.File)), r.FileName})
	case TypeFailed:
		return fmt.Sprintf("%+v", struct {
			Type       Type
			FailedCode int
		}{r.Type, r.FailedCode})
	case TypeRedirect:
		return fmt.Sprintf("%+v", struct {
			Type         Type
			RedirectURL  string
			RedirectCode int
		}{r.Type, r.RedirectURL, r.RedirectCode})
	default:
		return "<unknown>"
	}
}

// convert json data in response to JSON string
func (r *Response) ToJSON() string {
	marshal := jsonu.Marshal(r.Json)
	return marshal
}
