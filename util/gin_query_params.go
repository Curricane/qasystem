package util

import (
	"fmt"
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 把从gin获取的uint64数据转为int64
func GetQueryInt64(ctx *gin.Context, key string) (v int64, err error) {
	idstr, ok := ctx.GetQuery(key)
	if !ok {
		logger.Error("invalid params, not found key:%s", key)
		err = fmt.Errorf("invalid params, not found key:%s", key)
		return
	}

	v, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		logger.Error("invalid params, strconv.ParseInt failed, err:%v, str:%v",
			err, idstr)
		return
	}

	return
}
