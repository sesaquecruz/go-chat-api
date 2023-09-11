package dto

import "github.com/gin-gonic/gin"

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func AbortWithHttpError(c *gin.Context, status int, err error) {
	er := HttpError{
		Code:    status,
		Message: err.Error(),
	}

	c.AbortWithStatusJSON(status, er)
}
