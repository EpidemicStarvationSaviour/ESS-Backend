package rider

import (
	"ess/define"
	"ess/model/rider"
	"ess/service/rider_service"

	"ess/utils/authUtils"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary Rider Start Get Order
// @Tags    rider
// @Produce json
// @Param data body rider.RiderStart true "rider's position"
// @Success 200 {string} string "'success'"
// @Router  /rider/start [post]
func RiderStartGetOrder(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	RiderId := policy.GetId()
	var RSP rider.RiderStart
	if err := c.ShouldBind(&RSP); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	rider_service.GetRiderAvailable(RiderId)
	rider_service.RefreshRiderPosition(RiderId, RSP.AddressLat, RSP.AddressLng)
	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary Rider Stop Get Order
// @Tags    rider
// @Produce json
// @Success 200 {string} string "'success'"
// @Router  /rider/stop [post]
func RiderStopGetOrder(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	RiderId := policy.GetId()
	rider_service.GetRiderNotavailable(RiderId)
	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary Rider Upload Address
// @Tags    rider
// @Produce json
// @Param data body rider.RiderUploadAddress true "rider's position"
// @Success 200 {string} string "'success'"
// @Router  /rider/pos [post]
func RiderUploadAddressPort(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	RiderId := policy.GetId()
	var RSP rider.RiderUploadAddress
	if err := c.ShouldBind(&RSP); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	rider_service.RefreshRiderPosition(RiderId, RSP.AddressLat, RSP.AddressLng)
	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary Rider Check Whether there is a new order
// @Tags    rider
// @Produce json
// @Success 200 {object} rider.RiderQueryNewOrdersResp
// @Router  /rider/query [get]
func RiderQueryNewOrder(c *gin.Context) {
	err, rider := rider_service.QueryAvailableOrder()
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&rider))
}

// @Summary Rider's Feedback To New Order
// @Tags    rider
// @Produce json
// @Param data body rider.RiderFeedbackToNewOrder true "rider's feedback"
// @Success 200 {string} string "'success'"
// @Router  /rider/feedback [post]
func RiderFeedbackNeworder(c *gin.Context) {
	// claim, _ := c.Get(define.ESSPOLICY)
	// policy, _ := claim.(authUtils.Policy)
	// RiderId := policy.GetId()
	var RSP rider.RiderFeedbackToNewOrder
	if err := c.ShouldBind(&RSP); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	// rider_service.RiderFeedbackToOrder(RiderId, RSP.YesOrNo) // fake api
	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary feedback to new order
// @Tags    rider
// @Produce json
// @Param data body rider.FeedbackToOrder true "Your feedback"
// @Success 200 {string} string "'success'"
// @Router  /rider/groupfd [post]
func OrderFeedback(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	Uid := policy.GetId()
	var RSP rider.FeedbackToOrder
	if err := c.ShouldBind(&RSP); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	rider_service.RefreshOrderStatus(Uid, RSP)
}
