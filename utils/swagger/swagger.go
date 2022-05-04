package swagger

import (
	"ess/define"
	"ess/docs"
)

func Setup() {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "ESS API Documentation"
	docs.SwaggerInfo.Description = "all the json responses are in the form of {Status: \"success\", Msg: \"\",  Data: response}"
	docs.SwaggerInfo.Version = define.ESSAPIVERSION
}
