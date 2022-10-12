package query_service

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gin-gonic/gin"
)

func CreateLogCSVFile(logs []map[string]string, keys []string) string {
	fileName := fmt.Sprintf("room_query_%d.csv", time.Now().Unix())
	file, _ := os.Create(fileName)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(keys)
	for _, item := range logs {
		itemArray := make([]string, len(keys))
		for ki, k := range keys {
			itemArray[ki] = item[k]
		}
		writer.Write(itemArray[:])
	}
	return fileName
}

func DownloadHandler(c *gin.Context) {
	var user_query = struct {
		From int64  `form:"from" binding:"required"`
		To   int64  `form:"to" binding:"required"`
		Uuid string `form:"uuid" binding:"required"`
		Suid string `form:"suid"`
	}{}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc", user_query.Uuid, user_query.Suid)
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

func CutomQueryLogHandler(c *gin.Context) {
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
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc limit %d,%d", user_query.Uuid, user_query.Suid, user_query.Page*user_query.PageSize, user_query.PageSize)
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
