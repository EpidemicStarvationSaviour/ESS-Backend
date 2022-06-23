package category_service

import (
	"ess/model/category"
	"ess/model/item"
	"ess/utils/db"
	"ess/utils/logging"
)

func QueryCategoryById(cid int) *category.Category {
	var cat category.Category
	resinfo := db.MysqlDB.Where(&category.Category{CategoryId: cid}).First(&cat)
	logging.InfoF("Find %d category with cid %d\n", resinfo.RowsAffected, cid)
	return &cat
}

func QueryAllCategory() (*[]category.Category, error) {
	category := []category.Category{}
	if err := db.MysqlDB.Find(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func QueryMyCategory(uid int) (*[]category.Category, error) {
	cat := []category.Category{}
	var ite []item.Item
	db.MysqlDB.Where(&item.Item{ItemUserId: uid}).Find(&ite)
	for _, cate := range ite {
		var icat category.Category
		db.MysqlDB.Where(&category.Category{CategoryId: cate.ItemCategoryId}).First(&icat)
		cat = append(cat, icat)
	}
	return &cat, nil
}
func QueryCategoryByTypeId(fatherid int) ([]category.Category, error) {
	categorys := []category.Category{}
	err := db.MysqlDB.Where(&category.Category{CategoryFatherId: fatherid}).Find(&categorys).Error //where是干嘛的
	return categorys, err
}

func ValidCategory(cat category.CategoryCreateReq) bool {
	if len(cat.CategoryName) == 0 || cat.CategoryFatherId == 0 || cat.CategoryPrice < 0 {
		return false
	}
	if len(cat.CategoryAvatar) == 0 {
		return false
	}
	return true
}

func CreateNewCategory(cate *category.Category) error {
	tx := db.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(cate).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(cate).Updates(cate).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func DeleteCategoryByID(cid int) error {
	return db.MysqlDB.Where(&category.Category{CategoryId: cid}).Delete(&category.Category{CategoryId: cid}).Error
}

func DeleteItemByCiD(cid int) error {
	var ite []item.Item
	if err := db.MysqlDB.Where(&item.Item{ItemCategoryId: cid}).Find(&ite).Error; err != nil { //数据库的返回值
		return nil

	}
	for _, cite := range ite {
		cite.ItemAmount = 0
		db.MysqlDB.Model(&cite).Where(&item.Item{ItemCategoryId: cid}).Update("item_amount", 0)
	}
	var err error
	return err
}

//
func QueryCategoryByCid(cid int) (*category.Category, error) {
	var cat category.Category
	if err := db.MysqlDB.Where(&category.Category{CategoryId: cid}).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func QueryItemByCid(cid int) (*[]item.Item, error) {
	ite := []item.Item{}
	if err := db.MysqlDB.Where(&item.Item{ItemCategoryId: cid}).Find(&ite).Error; err != nil { //数据库的返回值
		return nil, err
	}
	return &ite, nil
}

func ModifyCategoryNumberByCid(uid int, cid int, number float64) int {
	var ite item.Item
	if err := db.MysqlDB.Where(&item.Item{ItemUserId: uid, ItemCategoryId: cid}).First(&ite).Error; err != nil {

		return 0
	}
	ite.ItemAmount = ite.ItemAmount + number

	if ite.ItemAmount < 0 {
		return 0
	}

	db.MysqlDB.Model(&ite).Where(&item.Item{ItemCategoryId: cid}).Update("item_amount", ite.ItemAmount)
	return 1
}

func GetCategoryTotalNumberByCid(cid int) float64 {
	var ite []item.Item
	var number float64
	if err := db.MysqlDB.Where(&item.Item{ItemCategoryId: cid}).Find(&ite).Error; err != nil {

		return 0
	}
	number = 0
	for _, item := range ite {
		number = number + item.ItemAmount
	}
	return number
}
func GetCategoryMyNumberByCid(cid int, uid int) float64 {
	var ite []item.Item
	var number float64
	if err := db.MysqlDB.Where(&item.Item{ItemCategoryId: cid, ItemUserId: uid}).Find(&ite).Error; err != nil {

		return 0
	}
	number = 0
	for _, item := range ite {
		number = number + item.ItemAmount
	}
	return number
}
