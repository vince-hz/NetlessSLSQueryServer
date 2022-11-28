package query_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gin-gonic/gin"
)

func DownloadHandler(c *gin.Context) {
	var user_query = struct {
		From         int64    `form:"from" binding:"required"`
		To           int64    `form:"to" binding:"required"`
		Uuid         string   `form:"uuid" binding:"required,len=32"`
		Keys         []string `form:"keys"  binding:"required"`
		Suid         string   `form:"suid"`
		FileType     string   `form:"fileType"`
		TimeLocation string   `form:"timeLocation"`
	}{
		FileType:     "csv",
		TimeLocation: "",
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
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc limit 10000", user_query.Uuid, suidQuery)
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
			"error": logError.Error(),
		})
		return
	} else {
		var filePath string
		if user_query.FileType == "csv" {
			filePath = CreateLogCSVFile(logResponse.Logs, user_query.Keys)
		} else if user_query.FileType == "xlsx" {
			filePath = CreateLogXLSXFile(logResponse.Logs, user_query.Keys, user_query.TimeLocation)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		}
		fileName := logFileName(user_query.Uuid, user_query.TimeLocation, user_query.FileType)
		c.FileAttachment(filePath, fileName)
		c.File(filePath)
		os.Remove(filePath)
	}
}

func logFileName(uuid string, timeLocation string, fileType string) string {
	var locationSuffix string
	if len(timeLocation) <= 0 {
		locationSuffix = "-UTC"
	} else {
		locationSuffix = "-timelocation:" + timeLocation
	}
	return fmt.Sprintf("%s%s.%s", uuid, locationSuffix, fileType)
}

type TimeCountDayItem struct {
	CumulativeSessionsCount int `json:"cumulativeSessionsCount"`
	Timestamp               int `json:"timestamp"`
}
type TimeCountResult struct {
	Msg []TimeCountDayItem `json:"msg"`
}

func TeamRooms(c *gin.Context) {
	var user_query = struct {
		From     int64  `form:"from" binding:"required"`
		To       int64  `form:"to" binding:"required"`
		Team     string `form:"team" binding:"required"`
		PageSize int    `form:"pageSize"`
		Page     int    `form:"page"`
		Region   string `form:"region"`
	}{
		PageSize: 30,
		Page:     1,
	}
	if err := c.ShouldBindQuery(&user_query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := fmt.Sprintf("team: %s| SELECT DISTINCT uuid from log limit %d,%d", user_query.Team, user_query.Page, user_query.PageSize)
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
	logResponse, logError, histogramResponse, _ := LogQuery(request)
	if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError.Error(),
		})
		return
	} else {
		var list []gin.H = make([]gin.H, logResponse.Count)
		for i := 0; i < int(logResponse.Count); i++ {
			uuid := logResponse.Logs[i]["uuid"]
			timeCountRequestPath := fmt.Sprintf("https://operation-server.netless.link/room/detail/day?region=%s&room=%s&token=%s", user_query.Region, uuid, env.WhiteOperationToken)
			response, _ := http.Get(timeCountRequestPath)
			body, _ := ioutil.ReadAll(response.Body)
			var res TimeCountResult
			json.Unmarshal(body, &res)
			var timeCount = 0
			var timestamp = 0
			if len(res.Msg) > 0 {
				for j := 0; j < len(res.Msg); j++ {
					timeCount += res.Msg[j].CumulativeSessionsCount
					if res.Msg[j].Timestamp > timestamp {
						timestamp = res.Msg[j].Timestamp
					}
				}
			}

			list[i] = gin.H{
				"uuid":      uuid,
				"timeCount": timeCount,
				"timestamp": timestamp,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"list":  list,
			"count": histogramResponse.Count,
		})
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
	hasSelect := strings.Contains(strings.ToLower(user_query.CustomQuery), "select")
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
			"error": logError.Error(),
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
		Uuid     string `form:"uuid" binding:"required,len=32"`
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
	pageStart := (user_query.Page - 1) * user_query.PageSize
	query := fmt.Sprintf("uuid: %s %s | SELECT * from log ORDER BY createdat asc limit %d,%d", user_query.Uuid, suidQuery, pageStart, user_query.PageSize)
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
			"error": hError.Error(),
		})
		return
	} else if logError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": logError.Error(),
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
		From         int64    `form:"from" binding:"required"`
		To           int64    `form:"to" binding:"required"`
		CustomQuery  string   `form:"customQuery" binding:"required"`
		Keys         []string `form:"keys"  binding:"required"`
		FileType     string   `form:"fileType"`
		TimeLocation string   `form:"timeLocation"`
	}{
		FileType:     "csv",
		TimeLocation: "",
	}
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
			"error": logError.Error(),
		})
		return
	} else {
		var filePath string
		if user_query.FileType == "csv" {
			filePath = CreateLogCSVFile(logResponse.Logs, user_query.Keys)
		} else if user_query.FileType == "xlsx" {
			filePath = CreateLogXLSXFile(logResponse.Logs, user_query.Keys, user_query.TimeLocation)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		}
		fileName := logFileName(user_query.CustomQuery, user_query.TimeLocation, user_query.FileType)
		c.FileAttachment(filePath, fileName)
		os.Remove(filePath)
	}
}
