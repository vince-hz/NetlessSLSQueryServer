package httpserver

import (
	"github.com/gin-gonic/gin"
)

func DefaultEngine() *gin.Engine {
	r := gin.Default()
	r.Use(func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		}
	}())
	return r
}
