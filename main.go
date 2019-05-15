package main

import (
    "github.com/gin-gonic/gin"
    "github.com/file-api-server/utils"
    "github.com/file-api-server/backends"
    "github.com/golang/glog"
    "net/http"
    "flag"
    "strconv"
)

var (
    config backends.Config
    user   *backends.S3UserInfo
    admin  *backends.S3UserInfo
)

type Bucket struct {
    BucketName  string `uri:"bucket" binding:"required"`
}

type File struct {
    BucketName  string `uri:"bucket" binding:"required"`
    FileName    string `uri:"file" binding:"required"`
}

func main() {
    //flag.Set("stderrthreshold", "INFO")
    flag.Parse()
    defer glog.Flush()

    utils.LoadConfig("config.yaml", &config)
    user = config.Users["user"]
    admin = config.Users["admin"]

    //gin.SetMode(gin.ReleaseMode)
    //r := gin.New()
    //r.Use(gin.Logger())
    //r.Use(gin.Recovery())
    r := gin.Default()
    r.GET("/buckets", listBuckets)
    r.POST("/bucket/:bucket", createBucket)
    r.GET("/bucket/:bucket", listBucket)
    r.DELETE("/bucket/:bucket", deleteBucket)
    r.GET("/policy/:bucket", getBucketPolicy)
    r.POST("file/:bucket", uploadFile)
    r.GET("file/:bucket/:file", downloadFile)
    r.DELETE("file/:bucket/:file", deleteFile)
    r.GET("/key/:bucket", getUserKey)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "ok",
        })
    })

    r.Run()
}

func listBuckets(c *gin.Context) {
    if buckets,err := user.ListBuckets(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to list buckets": err.Error()})
    } else {
        c.JSON(http.StatusOK, buckets)
    }
}

func createBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    if err := user.CreateBucket(c.Param("bucket")); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to create bucket": err.Error()})
    } else {
        c.String(http.StatusCreated, "")
    }
}

func deleteBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    if err := user.RemoveBucket(c.Param("bucket")); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to remove bucket": err.Error()})
    } else {
        c.String(http.StatusOK, "")
    }
}

func listBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    recursive,_ := strconv.ParseBool(c.DefaultQuery("recursive", "false"))
    if results,err := user.ListBucket(c.Param("bucket"), c.Query("prifix"), recursive); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to list files": err.Error()})
    } else {
        c.JSON(http.StatusOK, results)
    }
}

func getBucketPolicy(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    if policy,err := user.GetBucketPolicy(bucket.BucketName); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to get bucket policy": err.Error()})
    } else {
        c.String(http.StatusOK, policy)
    }
}

func uploadFile(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Bad Uri": err.Error()})
        return
    }

    file, err := c.FormFile("file");
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Bad Form Field": err.Error()})
        return
    }

    //if err := c.SaveUploadedFile(file, filename); err != nil {
    //	c.JSON(http.StatusBadRequest, gin.H{"msg": err})
    //	return
    //}

    if err := user.UploadFile(bucket.BucketName, file); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
    } else {
        c.String(http.StatusCreated, "")
    }
}

func downloadFile(c *gin.Context) {
    var f File

    if err := c.ShouldBindUri(&f); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }

    if object,err := user.DownloadFile(f.BucketName, f.FileName); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
    } else {
        if stat,err := object.Stat(); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        } else {
            //extraHeaders := map[string]string {
            //    "Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, stat.Key),
            //}
            c.DataFromReader(http.StatusOK, stat.Size, stat.ContentType, object, map[string]string{})
        }
    }
}

func deleteFile(c *gin.Context) {
    var f File

    if err := c.ShouldBindUri(&f); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }

    if err := user.RemoveFile(f.BucketName, f.FileName); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
    } else {
        c.String(http.StatusOK, "")
    }
}

func getUserKey(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }
    if key,err := admin.GetUserOfBucket(bucket.BucketName); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"failed to get key": err.Error()})
    } else {
        c.JSON(http.StatusOK, key)
    }
}
