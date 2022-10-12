package query_service

import sls "github.com/aliyun/aliyun-log-go-sdk"

var (
	ProjectName     = ""
	Endpoint        = ""
	LogStoreName    = ""
	AccessKeyID     = ""
	AccessKeySecret = ""
	Client          sls.ClientInterface
)
