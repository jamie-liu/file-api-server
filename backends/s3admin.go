package backends

import (
    "github.com/QuentinPerez/go-radosgw/pkg/api"
    "github.com/golang/glog"
    "fmt"
)

func initAdminOpsAPI(endpoint, accessKey, secretKey string) (*radosAPI.API, error) {
    return radosAPI.New(endpoint, accessKey, secretKey)
}

func (s3client *S3UserInfo) GetBucketStat(bucketName string) (*radosAPI.Stats, error) {
    api,err := initAdminOpsAPI(s3client.Endpoint, s3client.AccessKey, s3client.SecretKey)
    if err != nil {
        glog.Errorln(err)
        return nil, err
    }

    buckets, err := api.GetBucket(radosAPI.BucketConfig{Bucket: bucketName})
    if err != nil {
        glog.Errorln(err)
        return nil, err
    }
    for _,bucket := range buckets {
        return bucket.Stats, nil
    }
    return nil, fmt.Errorf("Bucket stats doesn't exist\n", bucketName)
}

func (s3client *S3UserInfo) GetUserOfBucket(bucketName string) (*S3Key, error) {
    api,err := initAdminOpsAPI(s3client.Endpoint, s3client.AccessKey, s3client.SecretKey)
    if err != nil {
        glog.Errorln(err)
        return nil, err
    }

    buckets, err := api.GetBucket(radosAPI.BucketConfig{Bucket: bucketName})
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
            return &S3Key{AccessKey: key.AccessKey, SecretKey: key.SecretKey, Endpoint: s3client.Endpoint}, nil
        }
    }
    return nil, fmt.Errorf("Bucket %s user key doesn't exist", bucketName)
}
