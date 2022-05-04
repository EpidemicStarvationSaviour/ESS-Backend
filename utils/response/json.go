package response

type JSONResponse struct {
	Status string      `json:"status"`
	Code   Code        `json:"-"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// data: body/nil
func JSONData(data interface{}) Response {
	return Response{
		Type: TypeJSON,
		Json: JSONResponse{
			Status: "success",
			Msg:    "",
			Data:   data,
		},
	}
}

// predefined error code
func JSONError(code Code) Response {
	return Response{
		Type: TypeJSON,
		Json: JSONResponse{
			Status: "error",
			Code:   code,
			Msg:    GetMsg(code),
			Data:   nil,
		},
	}
}

// customized error code
func JSONErrorWithMsg(msg string) Response {
	return Response{
		Type: TypeJSON,
		Json: JSONResponse{
			Status: "error",
			Code:   ERROR_DEFAULT,
			Msg:    msg,
			Data:   nil,
		},
	}
}
