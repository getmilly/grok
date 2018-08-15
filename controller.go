package api

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

//Controller ...
type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

//ResolveError ...
func ResolveError(context *gin.Context, err error) {
	if reflect.TypeOf(err) == reflect.TypeOf(Error{}) {
		status := http.StatusBadRequest
		if err.(Error).HTTPStatusCode != 0 {
			status = err.(Error).HTTPStatusCode
		}
		context.JSON(status, err)
		return
	}

	context.Status(http.StatusInternalServerError)
}

//BindingError ...
func BindingError(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, Error{Code: "000", Message: err.Error()})
}
