package query_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

type Env = struct {
	ProjectName     string `form:"ProjectName"`
	Endpoint        string
	LogStoreName    string
	AccessKeyID     string
	AccessKeySecret string
}

var (
	env    = Env{}
	Client sls.ClientInterface
)

func init() {
	file, _ := os.Open("./env.json")
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &env)
	Client = sls.CreateNormalInterface(env.Endpoint, env.AccessKeyID, env.AccessKeySecret, "")
}

func LogQuery(request sls.GetLogRequest) (*sls.GetLogsResponse, error, *sls.GetHistogramsResponse, error) {
	if env.ProjectName == "" {
		return nil, fmt.Errorf("ProjectName is empty"), nil, fmt.Errorf("ProjectName is empty")
	}
	h, hError := Client.GetHistograms(env.ProjectName, env.LogStoreName, request.Topic, request.From, request.To, request.Query)
	l, lError := Client.GetLogsV2(env.ProjectName, env.LogStoreName, &request)
	return l, lError, h, hError
}
