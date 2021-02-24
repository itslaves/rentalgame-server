package errors

import "github.com/gin-gonic/gin"

type ErrorCode int32

const (
	LoginRequired ErrorCode = 5400
	JoinRequired  ErrorCode = 5401
)

var errorMessage = map[ErrorCode]string{
	LoginRequired: "Login Required.",
	JoinRequired:  "Unauthorized member.",
}

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

func ErrorResponse(errorCode ErrorCode) gin.H {
	return ErrorResponseWithDetail(errorCode, nil)
}

func ErrorResponseWithDetail(errorCode ErrorCode, detail interface{}) gin.H {
	return gin.H{
		"error": Error{
			Code:    errorCode,
			Message: errorMessage[errorCode],
			Detail:  detail,
		},
	}
}
