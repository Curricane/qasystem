package db

import (
	"github.com/Curricane/logger"
	"qasystem/common"
)

func CreateFavoriteDir(fdir *common.FavoriteDir) (err error) {
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("create favorite dir failed, favorite dir:%#v, err:%v", fdir, err)
		return
	}

	//先查询相同的dir_name是否存在
	var dirCount int64
	sqlstr := `select count(dir_id) from favorite_dir where user_id=? and dir_name=?`
	err = tx.Get(&dirCount, sqlstr, fdir.UserId, fdir.DirName)
	if err != nil {
		logger.Error("select dir_name failed, err:%v, favoriteDir:%#v", err, fdir)
		return
	}

	if dirCount > 0 {
		tx.Rollback()
		err = ErrRecordExists
		return
	}

	sqlstr = `	insert into favorite_dir (user_id, dir_id, dir_name)
				values (?, ?, ?)`
	_, err = tx.Exec(sqlstr, fdir.UserId, fdir.DirId, fdir.DirName)
	if err != nil {
		logger.Error("insert favorite_dir failed, favorite_dir:%#v err:%v", fdir, err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Error("insert favorite_dir failed, favorite_dir:%#v err:%v", fdir, err)
		tx.Rollback()
		return
	}
	return
}

func CreateFavorite(f *common.Favorite) (err error) {
	tx, err := DB.Beginx()
	if err != nil {
		logger.Error("create favorite dir failed, favorite :%#v, err:%v", f, err)
		return
	}

	// 先查询相同的anser_id是否存在
	var fcount int64
	sqlstr := `select count(answer_id) from favorite where user_id=? and dir_id=?`
	err = tx.Get(&fcount, sqlstr, f.UserId, f.DirId)
	if err != nil {
		logger.Error("select dir_name failed, err:%v, favorite:%#v", err, f)
		return
	}

	if fcount > 0 {
		tx.Rollback()
		err = ErrRecordExists
		return
	}

	// 插入
	sqlstr = `	insert into favorite 
				  (user_id, dir_id, answer_id)
				values (?, ?, ?)`
	_, err = tx.Exec(sqlstr, f.UserId, f.DirId, f.AnswerId)
	if err != nil {
		logger.Error("insert favorite failed, favorite:%#v err:%v", f, err)
		tx.Rollback()
		return
	}

	// 收藏夹收藏数+1
	sqlstr = `update favorite_dir set count = count + 1 where dir_id = ?`
	_, err = tx.Exec(sqlstr, f.DirId)
	if err != nil {
		logger.Error("insert favorite failed, favorite:%#v err:%v", f, err)
		tx.Rollback()
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("insert favorite failed, favorite:%#v err:%v", f, err)
		tx.Rollback()
		return
	}
	return
}

func GetFavoriteDirList(userId int64) (favoriteDirList []*common.FavoriteDir, err error) {

	sqlstr := `	select dir_id, dir_name, count
	 			from favorite_dir
	  			where user_id=?`
	err = DB.Select(&favoriteDirList, sqlstr, userId)
	if err != nil {
		logger.Error("select favorite dir failed, err:%v", err)
		return
	}

	return
}

func GetFavoriteList(userId, dirId, offset, limit int64) (favoriteList []*common.Favorite, err error) {

	sqlstr := `	select dir_id, user_id, answer_id
				from favorite
				where user_id=? and dir_id=? limit ?, ?`
	err = DB.Select(&favoriteList, sqlstr, userId, dirId, offset, limit)
	if err != nil {
		logger.Error("select favorite list failed, err:%v", err)
		return
	}

	return
}
