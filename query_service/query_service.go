package query_service

import (
	"net/http"
	"os"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gin-gonic/gin"
)

func CustomQueryDownloadLog(c *gin.Context) {
	var user_query = struct {
		From        int64  `form:"from" binding:"required"`
		To          int64  `form:"to" binding:"required"`
		CustomQuery string `form:"customQuery" binding:"required"`
	}{}
	request := sls.GetLogRequest{
		From:     user_query.From,
		To:       user_query.To,
		Topic:    "",
		Lines:    0,
		Offset:   0,
		Reverse:  false,
		Query:    user_query.CustomQuery,
		PowerSQL: false,
	}
	keys := c.QueryArray("keys")
	logResponse, logError, _, _ := LogQuery(request)
	if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError,
		})
		return
	} else {
		fileName := CreateLogCSVFile(logResponse.Logs, keys)
		c.File(fileName)
		os.Remove(fileName)
	}
}
