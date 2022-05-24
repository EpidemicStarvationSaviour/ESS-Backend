package category_service

import (
	"ess/model/category"
	"ess/utils/db"
	"ess/utils/logging"
)

func QueryCategoryById(cid int) *category.Category {
	var cat *category.Category
	resinfo := db.MysqlDB.Where(&category.Category{CategoryId: cid}).First(&cat)
	logging.InfoF("Find %d category with cid %d\n", resinfo.RowsAffected, cid)
	return cat
}

// func AddCategoryGroupRelation(gid int, cid int) error {
// 	logging.InfoF("gid= %d, cid= %d\n", gid, cid)
// 	err := db.MysqlDB.Raw("INSERT INTO `group_category` () VALUES (?, ?)", gid, cid).Error
// 	return err
// }
