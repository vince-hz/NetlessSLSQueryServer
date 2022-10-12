package query_service

import (
	sls "github.com/aliyun/aliyun-log-go-sdk"
)

func LogQuery(request sls.GetLogRequest) (*sls.GetLogsResponse, error, *sls.GetHistogramsResponse, error) {
	Client = sls.CreateNormalInterface(Endpoint, AccessKeyID, AccessKeySecret, "")
	h, hError := Client.GetHistograms(ProjectName, LogStoreName, request.Topic, request.From, request.To, request.Query)
	l, lError := Client.GetLogsV2(ProjectName, LogStoreName, &request)
	return l, lError, h, hError
}
