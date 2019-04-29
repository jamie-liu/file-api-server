package utils

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/minio/minio-go"
)

func (s3client *S3Backend) CreateBucket(bucketName string) error {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	if err := conn.MakeBucket(bucketName, s3client.Location); err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (s3client *S3Backend) ListBuckets() ([]minio.BucketInfo, error) {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	buckets, err := conn.ListBuckets()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	return buckets, nil
}

func (s3client *S3Backend) RemoveBucket(bucketName string) error {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	if err := conn.RemoveBucket(bucketName); err != nil {
		glog.Errorln(err)
		return err
	}
	return nil
}

func (s3client *S3Backend) ListBucket(bucketName, prefix string) ([]FileStat, error) {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	if exists, err := conn.BucketExists(bucketName); err != nil {
		glog.Errorln(err)
		return nil, err
	} else if exists {
		// Create a done channel to control 'ListObjects' go routine.
		doneCh := make(chan struct{})

		// Indicate to our routine to exit cleanly upon return.
		defer close(doneCh)

		isRecursive := true
		objectCh := conn.ListObjects(bucketName, prefix, isRecursive, doneCh)
		results := []FileStat{}
		for object := range objectCh {
			results = append(results, FileStat{
				Name: object.Key,
				Size: object.Size,
				LastModified: object.LastModified,
			})
		}
		return results, nil
	}

	return nil, fmt.Errorf("Bucket doesn't exist: "+bucketName)
}

func (s3client *S3Backend) GetBucketPolicy(bucketName string) (string, error) {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		glog.Errorln(err)
		return "", err
	}

	bucketPolicy, err := conn.GetBucketPolicy(bucketName)
	if err != nil {
		glog.Errorln(err)
		return "", err
	}
	return bucketPolicy, nil
}