package controller

import (
	"github.com/gin-gonic/gin"
	"trama/internal/presentation/api/errorhandler"
)

func query(gc *gin.Context, target any) bool {
	if err := gc.ShouldBindQuery(target); err != nil {
		fail(gc, err)
		return false
	}
	return true
}

func path(gc *gin.Context, target any) bool {
	if err := gc.ShouldBindUri(target); err != nil {
		fail(gc, err)
		return false
	}
	return true
}

func body(gc *gin.Context, target any) bool {
	if err := gc.ShouldBindJSON(target); err != nil {
		fail(gc, err)
		return false
	}
	return true
}

func fail(gc *gin.Context, err error) {
	status, body := errorhandler.GlobalErrorHandler(err)
	gc.JSON(status, body)
}
