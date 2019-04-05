package api

import (
	"net/http"
	"reflect"

	"github.com/getmilly/grok/models"
	"github.com/gin-gonic/gin"
)

//Controller ...
type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

//ResolveError ...
func ResolveError(context *gin.Context, err error) {
	context.Error(err)

	if reflect.TypeOf(err) != reflect.TypeOf(models.Error{}) {
		context.Status(http.StatusInternalServerError)
		return
	}

	status := http.StatusBadRequest
	message := err.(models.Error)

	if message.HTTPStatusCode != 0 {
		status = message.HTTPStatusCode
	}

	context.JSON(status, message)
}

//BindingError ...
func BindingError(context *gin.Context, err error) {
	context.Error(err)
	message := models.Error{Code: "000", Message: err.Error()}
	context.JSON(http.StatusBadRequest, message)
}
