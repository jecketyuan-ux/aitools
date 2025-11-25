package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/playedu/playedu-go/pkg/constants"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginationResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    constants.SuccessCode,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    constants.SuccessCode,
		Message: message,
		Data:    data,
	})
}

func SuccessPagination(c *gin.Context, items interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    constants.SuccessCode,
		Message: "success",
		Data: PaginationResponse{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ParamError(c *gin.Context, message string) {
	Error(c, constants.ParamErrorCode, message)
}

func Unauthorized(c *gin.Context) {
	Error(c, constants.UnauthorizedCode, "unauthorized")
}

func PermissionDenied(c *gin.Context) {
	Error(c, constants.PermissionDeniedCode, "permission denied")
}

func NotFound(c *gin.Context, message string) {
	Error(c, constants.NotFoundCode, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, constants.ErrorCode, message)
}

func RateLimitExceeded(c *gin.Context) {
	Error(c, constants.RateLimitCode, "rate limit exceeded")
}
