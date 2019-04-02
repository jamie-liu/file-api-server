package main

import (
    "github.com/minio/minio-go"
    "github.com/QuentinPerez/go-radosgw/pkg/api"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    endpoint := "127.0.0.1:9000"
    accessKeyID := "test"
    secretAccessKey := "Test2017"
    useSSL := false

    minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
    if err != nil {
        log.Fatalln(err)
    }

    bucketName := "mymusic"
    location := "us-east-1"

    err = minioClient.MakeBucket(bucketName, location)

    if err != nil {
        exists, err := minioClient.BucketExists(bucketName)
        if err == nil && exists {
            log.Printf("We already own %s\n", bucketName)
        } else {
            log.Printf("The bucket name %s is not unique\n", bucketName)
            log.Fatalln(err)
        }
    }
    log.Printf("Successfully created %s\n", bucketName)

    objectName := "main.go"
    filePath := "./main.go"
    contentType := "application/text"

    n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType:contentType})
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

    api,_ := radosAPI.New(endpoint, accessKeyID, secretAccessKey)
    log.Printf("%#v\n", api)

    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
