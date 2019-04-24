package utils


import (
	"github.com/QuentinPerez/go-radosgw/pkg/api"
	"log"
)

func initAdminOpsAPI(endpoint, accessKeyID, secretAccessKey string) (*radosAPI.API, error) {
	api,err := radosAPI.New(endpoint, accessKeyID, secretAccessKey)
	log.Printf("%#v\n", api)
	return api,err
}