package category

import (
	"ess/define"
	"ess/model/address"
	"ess/model/category"
	"ess/model/item"

	"ess/service/address_service"
	"ess/service/category_service"
	"ess/service/user_service"

	"ess/utils/authUtils"
	//"ess/utils/logging"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @Summary get category list and total number
// @Tags    commodity
// @Produce json
// @Success 200 {object} category.CategoryInfoResp
// @Router  /commodity/list [get]
func GetCategoryList(c *gin.Context) {
	//claim, _ := c.Get(define.ESSPOLICY)
	//policy, _ := claim.(authUtils.Policy)
	var cate *[]category.Category
	var err error

	cate, err = category_service.QueryAllCategory()

	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var result category.CategoryAllResp
	result.CategoryList = make([]category.CategoryInfoResp, 0)

	for _, cat := range *cate {
		var data category.CategoryInfoResp
		data.Categorychild = make([]category.CategoryChildren, 0)
		data.CategoryNumber = 0
		if cat.CategoryLevel == 0 {
			copier.Copy(&data, &cat)
			var catechild []category.Category
			catechild, err = category_service.QueryCategoryByTypeId(cat.CategoryId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
				c.Abort()
				return
			}
			data.CategoryLevel = cat.CategoryId
			for _, categchild := range catechild {
				data.CategoryNumber++
				var catchild category.CategoryChildren

				copier.Copy(&catchild, &categchild)
				catchild.CategoryTotal = category_service.GetCategoryTotalNumberByCid(categchild.CategoryId)
				catchild.CategoryAvatar = categchild.CategoryImageUrl
				data.Categorychild = append(data.Categorychild, catchild)
			}
			result.CategoryList = append(result.CategoryList, data)
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result.CategoryList))
}

// @Summary add category
// @Tags    commodity
// @Produce json
// @Param data body category.CategoryCreateReq true "New Category Info"
// @Success 200 {object} category.CategoryCreateResp
// @Router  /commodity/add [post]
func CreateCate(c *gin.Context) {
	//claim, _ := c.Get(define.ESSPOLICY)
	//policy, _ := claim.(authUtils.Policy)
	var cate category.CategoryCreateReq
	if err := c.ShouldBind(&cate); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}
	valid := category_service.ValidCategory(cate)
	if !valid {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_VALID_USER_PARAM))
		c.Abort()
	}

	var cat category.Category
	cat.CategoryImageUrl = cate.CategoryAvatar
	cat.CategoryLevel = 1
	copier.Copy(&cat, &cate)
	if err := category_service.CreateNewCategory(&cat); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg(err.Error()))
		c.Abort()
	}

	resp := category.CategoryCreateResp{CategoryId: cat.CategoryId}
	c.Set(define.ESSRESPONSE, response.JSONData(resp))
}

// @Summary delete category
// @Tags    commodity
// @Produce json
// @Param data body category.CategoryDeleted true "Deleted Category id"
// @Success 200 {string} string "'success'"
// @Router  /commodity/delete [delete]
func DeleteCate(c *gin.Context) {
	var cid category.CategoryDeleted
	if err := c.ShouldBind(&cid); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}

	if err := category_service.DeleteCategoryByID(cid.CategoryId); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}
	// if err := category_service.DeleteItemByCiD(cid.CategoryId); err != nil {
	// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
	// 	c.Abort()
	// }
	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary get certain category details
// @Tags    commodity
// @Produce json
// @Param id path int true "edit category id"
// @Success 200 {object} category.CategoryCertainInfoResp
// @Router  /commodity/details/{id} [get]
func GetCateDetail(c *gin.Context) {
	//claim, _ := c.Get(define.ESSPOLICY)
	//policy, _ := claim.(authUtils.Policy)
	var cat category.CategoryCertainInfoResp
	var req_uri category.CategoryCertainInfoRep
	var err error

	var cate *category.Category
	if err := c.ShouldBindUri(&req_uri); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	cat.CategoryId = req_uri.CategoryId
	cate, err = category_service.QueryCategoryByCid(cat.CategoryId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	cat.CategoryFatherId = cate.CategoryFatherId
	cat.CategoryAvatar = cate.CategoryImageUrl
	cat.CategoryName = cate.CategoryName
	cat.CategoryPrice = cate.CategoryPrice

	var item *[]item.Item
	item, err = category_service.QueryItemByCid(cat.CategoryId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	var data category.CategoryCertainInfoResp
	data.CategoryDetails = make([]category.CategoryDetailsInfo, 0)
	copier.Copy(&data, &cat)
	var totalnumber float64
	totalnumber = 0
	for _, catitem := range *item {

		var catd category.CategoryDetailsInfo
		catd.StoreId = catitem.ItemUserId
		catd.Categorynumber = catitem.ItemAmount
		totalnumber = totalnumber + catd.Categorynumber
		var addr address.Address
		addr, err = address_service.QueryDefaultAddressByUserId(catd.StoreId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
			c.Abort()
			return
		}
		copier.Copy(&catd.CategoryAddress, &addr)
		catd.CategoryAddress = addr.AddressDetail
		catd.CategoryLat = addr.AddressLat
		catd.CategoryLng = addr.AddressLat
		usr := user_service.QueryUserById(catd.StoreId)
		catd.CategoryStorephone = usr.UserPhone

		data.CategoryDetails = append(data.CategoryDetails, catd)

	}

	data.CategoryTotal = totalnumber

	c.Set(define.ESSRESPONSE, response.JSONData(&data))

}

// @Summary get  my category info
// @Tags    commodity
// @Produce json
// @Success 200 {object} []category.CategoryMyInfoResp
// @Router  /commodity/my [get]
func GetMyCategoryDetails(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	userID := policy.GetId()
	var cate *[]category.Category
	var err error

	cate, err = category_service.QueryAllCategory()

	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	result := []category.CategoryMyInfoResp{}

	for _, cat := range *cate {
		var data category.CategoryMyInfoResp
		data.CategoryNumber = 0
		data.Categorychild = make([]category.CategoryMyChildren, 0)
		if cat.CategoryLevel == 0 {
			copier.Copy(&data, &cat)
			var catechild []category.Category
			catechild, err = category_service.QueryCategoryByTypeId(cat.CategoryId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
				c.Abort()
				return
			}
			data.CategoryLevel = cat.CategoryId
			for _, categchild := range catechild {
				var catchild category.CategoryMyChildren
				catchild.CategoryTotal = category_service.GetCategoryMyNumberByCid(categchild.CategoryId, userID)
				if catchild.CategoryTotal > 0 {
					data.CategoryNumber++

					copier.Copy(&catchild, &categchild)

					catchild.CategoryAvatar = categchild.CategoryImageUrl
					data.Categorychild = append(data.Categorychild, catchild)
				}
			}
			if data.CategoryNumber > 0 {
				result = append(result, data)
			}
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

// @Summary modify  my category info
// @Tags    commodity
// @Produce json
// @Param data body category.CategoryModifyRep true "Modify Category id"
// @Success 200
// @Router  /commodity/restock [post]
func ModifyCategoryNumber(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	userID := policy.GetId()
	var modcid category.CategoryModifyRep
	if err := c.ShouldBind(&modcid); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	yesorno := category_service.ModifyCategoryNumberByCid(userID, modcid.CategoryId, modcid.Categorynumber)
	if yesorno == 0 {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}
	if yesorno == 1 {

		c.Set(define.ESSRESPONSE, response.JSONData("success"))
	}
}
