package utils

import (
	"github.com/minio/minio-go"
)

func initConnection(endpoint, accessKeyID, secretAccessKey string, secure bool) (*minio.Client, error) {
	conn, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func initConnectionWithRegion(endpoint, accessKeyID, secretAccessKey string, secure bool, location string) (*minio.Client, error) {
	conn, err := minio.NewWithRegion(endpoint, accessKeyID, secretAccessKey, secure, location)
	if err != nil {
		return nil, err
	}
	return conn, nil
}