package utils


import (
	"github.com/QuentinPerez/go-radosgw/pkg/api"
	"github.com/golang/glog"
	"fmt"
)

func initAdminOpsAPI(endpoint, accessKeyID, secretAccessKey string) (*radosAPI.API, error) {
	return radosAPI.New(endpoint, accessKeyID, secretAccessKey)
}

func (s3client *S3Backend) GetBucketStat(bucketName string) {
	api,err := initAdminOpsAPI(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey)
	if err != nil {
		glog.Errorln(err)
	}

	config := radosAPI.BucketConfig{
		Bucket: bucketName,
		UID: "",
		Stats: true,
	}
	buckets, err := api.GetBucket(config)
	if err != nil {
		glog.Errorln(err)
	}
	for _,bucket := range buckets {
		glog.Errorf(bucket.Stats.Owner)
	}
}

func (s3client *S3Backend) GetUserFromBucket(bucketName string) (*Key, error) {
	api,err := initAdminOpsAPI(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	config := radosAPI.BucketConfig{
		Bucket: bucketName,
		UID: "",
		Stats: true,
	}
	buckets, err := api.GetBucket(config)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	for _,bucket := range buckets {
		user,err := api.GetUser(bucket.Stats.Owner)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		for _,key := range user.Keys {
			return &Key{AccessKey: key.AccessKey, SecretKey: key.SecretKey, Endpoint: key.User}, nil
		}
	}
	return nil, fmt.Errorf("Bucket user key doesn't exist\n", bucketName)
}
