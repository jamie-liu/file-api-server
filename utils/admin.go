package utils


import (
	"github.com/QuentinPerez/go-radosgw/pkg/api"
	"github.com/golang/glog"
)

func initAdminOpsAPI(endpoint, accessKeyID, secretAccessKey string) (*radosAPI.API, error) {
	api,err := radosAPI.New(endpoint, accessKeyID, secretAccessKey)
	glog.Infof("%#v\n", api)
	return api,err
}