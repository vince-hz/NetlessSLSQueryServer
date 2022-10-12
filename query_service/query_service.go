package query_service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gin-gonic/gin"
)

func DownloadHandler(c *gin.Context) {
	var user_query = struct {
		From int64    `form:"from" binding:"required"`
		To   int64    `form:"to" binding:"required"`
		Uuid string   `form:"uuid" binding:"required"`
		Keys []string `form:"keys"  binding:"required"`
		Suid string   `form:"suid"`
	}{}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var suidQuery string
	if len(user_query.Suid) > 0 {
		suidQuery = fmt.Sprintf("and suid: %s", user_query.Suid)
	} else {
		suidQuery = ""
	}
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc", user_query.Uuid, suidQuery)
	request := sls.GetLogRequest{
		From:     user_query.From,
		To:       user_query.To,
		Topic:    "",
		Lines:    0,
		Offset:   0,
		Reverse:  false,
		Query:    query,
		PowerSQL: false,
	}

	logResponse, logError, _, _ := LogQuery(request)
	if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError,
		})
		return
	} else {
		fileName := CreateLogCSVFile(logResponse.Logs, user_query.Keys)
		c.File(fileName)
		os.Remove(fileName)
	}
}

func CustomQueryLogHandler(c *gin.Context) {
	var user_query = struct {
		From        int64  `form:"from" binding:"required"`
		To          int64  `form:"to" binding:"required"`
		CustomQuery string `form:"customQuery" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hasSelect := strings.Contains(user_query.CustomQuery, "select")
	if !hasSelect {
		c.JSON(http.StatusBadRequest, gin.H{"error": "custom query should conatins select statement"})
		return
	}
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
	logResponse, logError, _, _ := LogQuery(request)
	if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"list":  logResponse.Logs,
			"count": logResponse.Count,
		})
	}
}

func LogHandler(c *gin.Context) {
	var user_query = struct {
		From     int64  `form:"from" binding:"required"`
		To       int64  `form:"to" binding:"required"`
		Uuid     string `form:"uuid" binding:"required"`
		Suid     string `form:"suid"`
		PageSize int    `form:"pageSize"`
		Page     int    `form:"page"`
	}{
		PageSize: 30,
		Page:     1,
	}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var suidQuery string
	if len(user_query.Suid) > 0 {
		suidQuery = fmt.Sprintf("and suid: %s", user_query.Suid)
	} else {
		suidQuery = ""
	}
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc limit %d,%d", user_query.Uuid, suidQuery, user_query.Page*user_query.PageSize, user_query.PageSize)
	request := sls.GetLogRequest{
		From:     user_query.From,
		To:       user_query.To,
		Topic:    "",
		Lines:    0,
		Offset:   0,
		Reverse:  false,
		Query:    query,
		PowerSQL: false,
	}

	logResponse, logError, histogramResponse, hError := LogQuery(request)
	if hError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": hError,
		})
		return
	} else if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"list":  logResponse.Logs,
			"count": histogramResponse.Count,
		})
	}
}

func CustomQueryDownloadLog(c *gin.Context) {
	var user_query = struct {
		From        int64    `form:"from" binding:"required"`
		To          int64    `form:"to" binding:"required"`
		CustomQuery string   `form:"customQuery" binding:"required"`
		Keys        []string `form:"keys"  binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	logResponse, logError, _, _ := LogQuery(request)
	if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError,
		})
		return
	} else {
		fileName := CreateLogCSVFile(logResponse.Logs, user_query.Keys)
		c.File(fileName)
		os.Remove(fileName)
	}
}
