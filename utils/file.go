package utils

import (
	"github.com/minio/minio-go"
	"log"
	"fmt"
	"mime/multipart"
	"path/filepath"
)

func (s3client *S3Backend) UploadLocalFile(bucketName, objectName, filePath string) error {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		log.Println(err)
		return err
	}

	if exists, err := conn.BucketExists(bucketName); err != nil {
		log.Println(err)
		return err
	} else if !exists {
		err = conn.MakeBucket(bucketName, s3client.Location)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if n, err := conn.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{}); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
		return nil
	}
}

func (s3client *S3Backend) UploadFile(bucketName string, file *multipart.FileHeader) error {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		log.Println(err)
		return err
	}

	if exists, err := conn.BucketExists(bucketName); err != nil {
		log.Println(err)
		return err
	} else if !exists {
		err = conn.MakeBucket(bucketName, s3client.Location)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	filename := filepath.Base(file.Filename)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if n, err := conn.PutObject(bucketName, filename, src ,file.Size, minio.PutObjectOptions{}); err != nil {
		log.Println(err)
		return err
	} else {
		log.Printf("Successfully uploaded %s of size %d\n", filename, n)
		return nil
	}
}

func (s3client *S3Backend) DownloadFile(bucketName, objectName string) (*minio.Object, error) {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if exists, err := conn.BucketExists(bucketName); err != nil {
		log.Println(err)
		return nil, err
	} else if !exists {
		return nil, fmt.Errorf("Bucket %s doesn't exist\n", bucketName)
	}

	//if err := conn.FGetObject(bucketName, objectName, filePath, minio.GetObjectOptions{}); err != nil {
	return conn.GetObject(bucketName, objectName, minio.GetObjectOptions{})
}

func (s3client *S3Backend) RemoveFile(bucketName, objectName string) error {
	conn, err := initConnectionWithRegion(s3client.Endpoint, s3client.AccessKeyID, s3client.SecretAccessKey, s3client.SSL, s3client.Location)
	if err != nil {
		log.Println(err)
		return err
	}

	if exists, err := conn.BucketExists(bucketName); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("Bucket %s doesn't exist\n", bucketName)
	}

	if err := conn.RemoveObject(bucketName, objectName); err != nil {
		log.Println(err)
		return err
	}
	return nil
}