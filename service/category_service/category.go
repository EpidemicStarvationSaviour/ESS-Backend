package category_service

import (
	"ess/model/category"
	"ess/utils/db"
	"ess/utils/logging"
)

func QueryCategoryById(cid int) *category.Category {
	var cat *category.Category
	resinfo := db.MysqlDB.Where("cid = ?", cid).First(&cat)
	logging.InfoF("Find %d category with cid %d\n", resinfo.RowsAffected, cid)
	return cat
}
