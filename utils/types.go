package utils

import "time"

type S3Backend struct{
	Endpoint string
	AccessKeyID string
	SecretAccessKey string
	SSL bool
	Location string
}

type FileStat struct{
	Name 			string 		`json:"name"`
	Size 			int64  		`json:"size"`
	LastModified 	time.Time  	`json:"lastModified"`
}