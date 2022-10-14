package main

import (
	"netless/slsquery/httpserver"
	"netless/slsquery/query_service"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	args := os.Args
	if len(args) >= 2 {
		if args[1] == "release" {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	router := httpserver.DefaultEngine()
	router.GET("/logs", query_service.LogHandler)
	router.GET("/customQuery/logs", query_service.CustomQueryLogHandler)
	router.GET("/downloadLogs", query_service.DownloadHandler)
	router.GET("/customQuery/downloadLogs", query_service.CustomQueryDownloadLog)
	router.Run(":8080")
}
