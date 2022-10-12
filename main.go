package main

import (
	"netless/slsquery/httpserver"
	"netless/slsquery/query_service"
)

func main() {
	router := httpserver.DefaultEngine()
	router.GET("/logs", query_service.LogHandler)
	router.GET("/customQuery/Logs", query_service.CustomQueryDownloadLog)
	router.GET("/downloadLogs", query_service.DownloadHandler)
	router.GET("/customQuery/downloadLogs", query_service.CustomQueryDownloadLog)
	router.Run(":8080")
}
