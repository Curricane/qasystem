package db

import (
	"database/sql"
	"github.com/Curricane/logger"
	"qasystem/common"
	"github.com/jmoiron/sqlx"
)

/*
GetCategoryList 获取分类列表
return categoryList type slice
*/
func GetCategoryList() (categoryList []*common.Category, err error) {
	sqlstr := "select category_id, category_name from category"
	err = DB.Select(&categoryList, sqlstr)
	if err == sql.ErrNoRows {
		logger.Warn("there is no category result in db")
		err = nil
		return
	}
	if err != nil {
		logger.Error("failed to '%s', err is: %#v", sqlstr, err)
		return
	}
	return
}

func MGetCategory(categoryIds []int64) (categoryMap map[int64]*common.Category, err error) {
	sqlstr := "select category_id, category_name from category where category_id in (?)"
	var interfaceSlice []interface{}
	for _, c := range categoryIds {
		interfaceSlice = append(interfaceSlice, c)
	}
	insqlStr, params, err := sqlx.In(sqlstr, interfaceSlice...)
	if err != nil {
		logger.Error("sqlx.in failed, sqlstr:%v, err:%v", sqlstr, err)
		return
	}
	categoryMap = make(map[int64]*common.Category, len(categoryIds))
	var categoryList []*common.Category
	err = DB.Select(&categoryList, insqlStr, params...)
	if err != nil {
		logger.Error("MGetCategory  failed, sqlstr:%v, category_ids:%v, err:%v",
			sqlstr, categoryIds, err)
		return
	}

	for _, v := range categoryList {
		categoryMap[v.CategoryId] = v
	}
	return
}
