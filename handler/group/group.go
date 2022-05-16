package group

import (
	"ess/define"
	"ess/model/group"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary get groups I joined
// @Tags	group
// @Produce json
// @Success 200 {object} group.GroupInfoResp
// @Router /group/list [get]
func GetMyGroup(c *gin.Context) {
	var groupcondition group.GroupInfoReq
	if err := c.ShouldBind(&groupcondition); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
}
