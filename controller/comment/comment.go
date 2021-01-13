package comment

import (
	"fmt"
	"github.com/Curricane/logger"
	"github.com/gin-gonic/gin"
	"html"
	"qasystem/common"
	"qasystem/dal/db"
	"qasystem/id_gen"
	"qasystem/middleware/account"
	"qasystem/util"
	"strconv"
	"strings"
)

const (
	MinCommentContentSize = 10
)

/*检查传入参数Comment

 */
func checkComment(cmt *common.Comment) (err error) {
	if cmt == nil {
		err = fmt.Errorf("comment is nil")
		return
	}

	if len(cmt.Content) <= MinCommentContentSize || cmt.QuestionId == 0 ||
		cmt.ReplyCommentId == 0 || cmt.ParentId == 0 {
		err = fmt.Errorf("comment is invalid, len:%v, qid:%v", len(cmt.Content), cmt.QuestionId)
		return
	}
	return
}

/* 发表level 2评论
type Comment struct {
	CommentId       int64     // 无，id_gen生成
	Content         string    // 前端传过来
	AuthorId        int64     // 无，从session中获取
	LikeCount       int       // 无，不需要
	CommentCount    int       // 无，不需要
	CreateTime      time.Time // 无，不需要
	ParentId        int64     // 前端传过来，level 2时 为被评论的comment_id
	QuestionId      int64     // 前端传过来，回答/问题的id
	ReplyAuthorId   int64     // 前端传过来，level2评论，被评论的评论者id，也可以根据ReplyCommentId查询
	ReplyCommentId  int64     // 前端传过来
	AuthorName      string    // 前端传过来
	ReplyAuthorName string    // 前端传过来
	QuestionIdStr   string    // 前端传过来
}
*/
func PostReplyHandle(ctx *gin.Context) {
	var cmt common.Comment
	err := ctx.BindJSON(&cmt)
	if err != nil {
		logger.Error("bind json failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}
	logger.Debug("bind json succ, comment:%#v", cmt)

	if err = checkComment(&cmt); err != nil {
		logger.Error("invalid comment, err is:%v", err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	// step1 补充结构数据
	// 1.1 获取AuthorId，即UserId
	userId, err := account.GetUserId(ctx)
	if err != nil || userId <= 0 {
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	cmt.AuthorId = userId

	// 1.2 针对content做一个转义，防止xss漏洞
	cmt.Content = html.EscapeString(cmt.Content)

	// 1.3 生成评论id
	cid, err := id_gen.GetId()
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("id_gen.GetId failed, comment:%#v, err:%v", cmt, err)
		return
	}
	cmt.CommentId = int64(cid)

	err = db.CreateReplyComment(&cmt)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("CreatePostComment failed, comment:%#v, err:%v", cmt, err)
		return
	}

	util.ResponseSuccess(ctx, nil)
}

/* 发表level 1评论
type Comment struct {
	CommentId       int64     // 无，id_gen生成
	Content         string    // 前端传过来
	AuthorId        int64     // 无，从session中获取
	LikeCount       int       // 无，不需要
	CommentCount    int       // 无，不需要
	CreateTime      time.Time // 无，不需要
	ParentId        int64     // 前端传过来，level 1时 为0
	QuestionId      int64     // 前端传过来
	ReplyAuthorId   int64     // 前端传过来，level1评论，ReplyAuthorId是问题/回答者id
	ReplyCommentId  int64     // 前端传过来
	AuthorName      string    // 前端传过来
	ReplyAuthorName string    // 前端传过来
	QuestionIdStr   string    // 前端传过来
}
 */
func PostCommentHandle(ctx *gin.Context) {
	var cmt common.Comment
	err := ctx.BindJSON(&cmt)
	if err != nil {
		logger.Error("bind json failed, err:%v", err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}
	logger.Debug("bind json succ, comment:%#v", cmt)

	// 获取问题/回答id
	cmt.QuestionId, err = strconv.ParseInt(cmt.QuestionIdStr, 10, 64)
	if err != nil {
		logger.Error("cmt.QuestionIdStr is:%v cannt convert to int, err is:%v",cmt.QuestionIdStr, err)
		util.ResponseError(ctx, util.ErrCodeParameter)
		return
	}

	// 获取AuthorId，即UserId
	userId, err := account.GetUserId(ctx)
	if err != nil || userId <= 0 {
		util.ResponseError(ctx, util.ErrCodeNotLogin)
		return
	}
	cmt.AuthorId = userId

	// 针对content做一个转义，防止xss漏洞
	cmt.Content = html.EscapeString(cmt.Content)

	// 生成评论id
	cid, err := id_gen.GetId()
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("id_gen.GetId failed, comment:%#v, err:%v", cmt, err)
		return
	}
	cmt.CommentId = int64(cid)

	err = db.CreatePostComment(&cmt)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("CreatePostComment failed, comment:%#v, err:%v", cmt, err)
		return
	}

	util.ResponseSuccess(ctx, nil)
}

func CommentListHandle(ctx *gin.Context) {

	if ctx == nil {
		logger.Error("ctx is nil")
		return
	}

	//解析answer_id
	answerIdStr, ok := ctx.GetQuery("answer_id")
	answerIdStr = strings.TrimSpace(answerIdStr)
	if ok == false || len(answerIdStr) == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid answer id, val:%v", answerIdStr)
		return
	}
	answerId, err := strconv.ParseInt(answerIdStr, 10, 64)
	if err != nil || answerId == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid answer id, val:%v", answerIdStr)
		return
	}
	logger.Debug("get query answer_id succ, val:%v", answerIdStr)

	//解析offset
	var offset int64
	offsetStr, ok := ctx.GetQuery("offset")
	offsetStr = strings.TrimSpace(offsetStr)
	if ok == false || len(offsetStr) == 0 {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	logger.Debug("get query offset succ, val:%v", offsetStr)

	//解析limit
	var limit int64
	limitStr, ok := ctx.GetQuery("limit")
	limitStr = strings.TrimSpace(limitStr)
	if ok == false || len(limitStr) == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}

	logger.Debug("get query limit succ, val:%v", limitStr)

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}
	logger.Debug("get query limit succ, val:%v", limitStr)

	//获取一级评论列表
	commentList, count, err := db.GetCommentList(answerId, offset, limit)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("GetCommentList failed, answer_id:%v  err:%v", answerId, err)
		return
	}

	var userIdList []int64
	for _, v := range commentList {
		userIdList = append(userIdList, v.AuthorId, v.ReplyAuthorId)
	}

	userList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("GetUserInfoList failed, answer_id:%v  err:%v", answerId, err)
		return
	}

	userInfoMap := make(map[int64]*common.UserInfo, len(userIdList))
	for _, user := range userList {
		userInfoMap[int64(user.UserId)] = user
	}

	for _, v := range commentList {
		user, ok := userInfoMap[v.AuthorId]
		if ok {
			v.AuthorName = user.Username
		}

		user, ok = userInfoMap[v.ReplyAuthorId]
		if ok {
			v.ReplyAuthorName = user.Username
		}
	}

	var apiCommentList = &common.ApiCommentList{}
	apiCommentList.Count = count
	apiCommentList.CommentList = commentList

	util.ResponseSuccess(ctx, apiCommentList)
}

// 根据前端传过来的comment_id offset limit，获取评论列表ApiComentList
func ReplyListHandle(ctx *gin.Context) {
	// 解析comment_id
	cmtIdStr, ok := ctx.GetQuery("comment_id")
	cmtIdStr = strings.TrimSpace(cmtIdStr)
	if !ok || len(cmtIdStr) == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid comment id, val:%v", cmtIdStr)
		return
	}
	cmtId, err := strconv.ParseInt(cmtIdStr, 10, 64)
	if err != nil || cmtId == 0 {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("valid comment id, val:%v", cmtId)
		return
	}
	logger.Debug("get query commentIdStr succ, val:%v", cmtIdStr)

	//解析offset
	var offset int64
	offsetStr, ok := ctx.GetQuery("offset")
	offsetStr = strings.TrimSpace(offsetStr)
	if ok == false || len(offsetStr) == 0 {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	offset, err = strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
		logger.Error("invalid offset, val:%v", offsetStr)
	}
	logger.Debug("get query offset succ, val:%v", offsetStr)

	//解析limit
	var limit int64
	limitStr, ok := ctx.GetQuery("limit")
	limitStr = strings.TrimSpace(limitStr)
	if ok == false || len(limitStr) == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}
	logger.Debug("get query limit succ, val:%v", limitStr)

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit == 0 {
		limit = 10
		logger.Error("valid limit, val:%v", limitStr)
	}
	logger.Debug("get query limit succ, val:%v", limitStr)

	// 获取回复列表
	cmtList, cnt, err := db.GetReplyList(cmtId, offset, limit)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("GetCommentList failed, commentId:%v  err:%v", cmtId, err)
		return
	}

	// 获取出现的用户信息
	var userIdList []int64
	for _, v := range cmtList {
		userIdList = append(userIdList, v.AuthorId, v.ReplyAuthorId)
	}
	userList, err := db.GetUserInfoList(userIdList)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("GetUserInfoList failed, answer_id:%v  err:%v", cmtId, err)
		return
	}
	userInfoMap := make(map[int64]*common.UserInfo, len(userIdList))
	for _, user := range userList {
		userInfoMap[int64(user.UserId)] = user
	}

	// 获取cmtList中 cmt的AuthorName ReplyAuthorName
	for _, v := range cmtList {
		user, ok := userInfoMap[v.AuthorId]
		if ok {
			v.AuthorName = user.Username
		}

		user, ok = userInfoMap[v.ReplyAuthorId]
		if ok {
			v.ReplyAuthorName = user.Username
		}
	}

	var apiCommentList = &common.ApiCommentList{}
	apiCommentList.Count = cnt
	apiCommentList.CommentList = cmtList

	util.ResponseSuccess(ctx, apiCommentList)
}

func LikeHandle(ctx *gin.Context) {
	var like common.Like
	err := ctx.BindJSON(&like)
	if err != nil {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("like handler failed, err:%v", err)
		return
	}

	if like.Id == 0 || (like.LikeType != common.LikeTypeAnswer && like.LikeType != common.LikeTypeComment) {
		util.ResponseError(ctx, util.ErrCodeParameter)
		logger.Error("invalid like paramter, data:%#v", like)
		return
	}

	switch like.LikeType {
	case common.LikeTypeAnswer:
		err = db.UpdateAnswerLikeCount(like.Id)
	case common.LikeTypeComment:
		err = db.UpdateCommentLikeCount(like.Id)
	default:
		logger.Error("cannt match LikeType")
	}

	if err != nil {
		util.ResponseError(ctx, util.ErrCodeServerBusy)
		logger.Error("update like count failed, err:%v, data:%#v", err, like)
		return
	}

	util.ResponseSuccess(ctx, nil)
}

