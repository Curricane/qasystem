package db

import (
	"database/sql"
	"fmt"
	"qasystem/common"
	"qasystem/util"
)

// 以后可以改成根据用户信息生成一个salt，并存到数据库中
var PasswordSalt string = "123456789"

func Login(user *common.UserInfo) (err error) {

	//先保存用户密码
	originPassword := user.Password

	sqlstr := "select username, password, user_id from user where username=?"
	err = DB.Get(user, sqlstr, user.Username)
	// 查询数据库出错
	if (err != nil ) && (err != sql.ErrNoRows) {
		fmt.Println("failed to get username, password, err is:", err)
		return
	}

	// 查询不到用户
	if err == sql.ErrNoRows {
		err = ErrUserNotExits
		return
	}

	// 用户密码+salt 取md5值
	passwd := originPassword + PasswordSalt
	originPasswordSalt := util.Md5([]byte(passwd))

	// 对比MD5值
	if originPasswordSalt != user.Password {
		err = ErrUserPasswordWrong
		return
	}

	return
}

func Register(user *common.UserInfo) (err error) {
	if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Email) == 0 {
		err = fmt.Errorf("invalid UserInfo")
		return
	}

	var count int64
	sqlstr := "select count(user_id) from user where username = ?"
	err = DB.Get(&count, sqlstr, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if count > 0 {
		err = ErrUserExits
		return
	}

	passwd := user.Password + PasswordSalt

	dbPassword := util.Md5([]byte(passwd))

	// 之后可以替换为orm方式
	sqlstr = "insert into user(username, password, email , user_id, sex, nickname) values(?,?,?,?,?,?)"
	_, err = DB.Exec(sqlstr, user.Username, dbPassword, user.Email, user.UserId, user.Sex, user.Nickname)
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("failed to insert userinfo to db")
		return err
	}

	return
}