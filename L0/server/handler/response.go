package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// As `gin.Context.AbortWithStatusJSON()` accepts an epty interface and
// waits for the JSON object, it is crucial to create a custom
// error type such as following
type custom_error struct {
	Message string `json:"message"`
}

// Write error as json with and abort the following operations
// 
// Accepts context, status code (200, 400, 500) and erros message
func throwError (c *gin.Context, code int, msg string) {
	// Writes the error into logrus (external logger)
	logrus.Error(msg)
	// Aborts following operations
	c.AbortWithStatusJSON(code, custom_error{msg})
}