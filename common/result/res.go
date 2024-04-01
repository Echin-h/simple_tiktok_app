package result

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Any        any    `json:"any"`
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// 这个是对返回任意数值的包装
type R map[string]any

var (
	// 200 OK
	SuccessStatus = newStatus(20000, "success")

	// 400 BAD
	ErrEmpty                  = newStatus(40000, "username or password is empty")
	QueryParamErrorStatus     = newStatus(40001, "请求的参数错误")
	LoginErrorStatus          = newStatus(40002, "登录发生错误")
	RegisterErrorStatus       = newStatus(40003, "注册发生错误")
	UsernameExitErrorStatus   = newStatus(40004, "用户名已存在")
	TokenErrorStatus          = newStatus(40005, "token 错误")
	InfoErrorStatus           = newStatus(40006, "无法获取该用户信息")
	FileErrorStatus           = newStatus(40007, "文件上传失败")
	PublishErrorStatus        = newStatus(40008, "发布时出现错误")
	FeedErrorStatus           = newStatus(40009, "获取视频流出错")
	EmptyErrorStatus          = newStatus(40010, "用户名或密码为空") // should be useless
	FollowErrorStatus         = newStatus(40011, "关注失败")
	FavoriteErrorStatus       = newStatus(40012, "点赞失败")
	FollowListErrorStatus     = newStatus(40013, "获取关注列表时发生了错误")
	PublishListErrorStatus    = newStatus(40014, "获取发布列表时发生了错误")
	CommentPublishErrorStatus = newStatus(40015, "发布评论时发生了错误")

	// 401 WITHOUT PERMISSION
	NoLoginErrorStatus = newStatus(40101, "用户未登录")

	// 403 ILLEGAL OPERATION
	PermissionErrorStatus = newStatus(40301, "操作非法")

	// 404 NOT FOUND
	CommentNotExitErrorStatus = newStatus(40401, "评论不存在")
	VideoNotExitErrorStatus   = newStatus(40402, "视频不存在")

	// 500 INTERNAL ERROR
	ServerErrorStatus = newStatus(50001, "服务器内部错误")
)

func (s Status) Code() int {
	return s.StatusCode
}

func (s Status) Msg() string {
	return s.Message
}

func (s Status) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", s.StatusCode, s.Message)
}

func newStatus(code int, msg string) Status {
	return Status{StatusCode: code, Message: msg}
}

func Resp(c *gin.Context, status Status, data ...any) {
	c.JSON(status.StatusCode/100, Response{status.Message, data})
	c.Abort()
}

// 判断错误是否相同
func Is(status Status, err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(Status); ok {
		return e.StatusCode == status.StatusCode
	}
	return false
}

// 自创的话
//func Success(c *gin.Context) {
//	Resp(c, Status{
//		StatusCode: http.StatusOK,
//		Message: "success",
//	}, "happy")
//}

// 平常使用
//func Fail(c *gin.Context, PermissionErrorStatus Status){
//	Resp(c, PermissionErrorStatus, "sad")
//}
