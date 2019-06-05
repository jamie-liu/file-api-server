package main

import (
    "github.com/gin-gonic/gin"
    "github.com/file-api-server/utils"
    "github.com/file-api-server/backends"
    "github.com/golang/glog"
    "github.com/rakyll/statik/fs"
    "net/http"
    "flag"
    "strconv"
    "os"
    _ "github.com/file-api-server/statik"
)

var (
    config backends.Config
    yamlConfigFile string
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

    if yamlConfigFile = os.Getenv("CONFIG"); yamlConfigFile == "" {
        glog.Fatalf("CONFIG env not set")
    }
    utils.LoadConfig(yamlConfigFile, &config)
    //utils.GenerateRSAKey(2048, "private.pem", "public.pem")
    //for _,v := range config.Users {
    //    v.AccessKey,_ = utils.RsaDecrypt(v.AccessKey, "private.pem")
    //    v.SecretKey,_ = utils.RsaDecrypt(v.SecretKey, "private.pem")
    //}

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
    r.POST("/user/:user", setUserKey)
    r.POST("/reload", func(c *gin.Context) {
        if err := utils.LoadConfig(yamlConfigFile, &config); err != nil {
            c.JSON(http.StatusNotFound, gin.H{"failed to reload config": err.Error()})
        } else {
            c.JSON(http.StatusOK, "")
        }
    })
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "health": "ok",
        })
    })

    //r.Static("/swagger-ui/", "./swagger-ui")
    if statikFS,err := fs.New(); err != nil {
        glog.Errorln(err)
    } else {
        r.StaticFS("/swagger-ui/", statikFS)
    }

    r.Run()
}

func listBuckets(c *gin.Context) {
    cloud := config.GetUserConfig("cloud")
    if buckets,err := cloud.ListBuckets(); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to get bucket list": err.Error()})
    } else {
        c.JSON(http.StatusOK, buckets)
    }
}

func createBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cloud := config.GetUserConfig("cloud")
    if err := cloud.CreateBucket(c.Param("bucket")); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to create bucket": err.Error()})
    } else {
        c.String(http.StatusCreated, "")
    }
}

func deleteBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cloud := config.GetUserConfig("cloud")
    if err := cloud.RemoveBucket(c.Param("bucket")); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to remove bucket": err.Error()})
    } else {
        c.String(http.StatusOK, "")
    }
}

func listBucket(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    recursive,_ := strconv.ParseBool(c.DefaultQuery("recursive", "false"))
    cloud := config.GetUserConfig("cloud")
    if results,err := cloud.ListBucket(c.Param("bucket"), c.Query("prifix"), recursive); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to list bucket": err.Error()})
    } else {
        c.JSON(http.StatusOK, results)
    }
}

func getBucketPolicy(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    cloud := config.GetUserConfig("cloud")
    if policy,err := cloud.GetBucketPolicy(bucket.BucketName); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to get bucket policy": err.Error()})
    } else {
        c.String(http.StatusOK, policy)
    }
}

func uploadFile(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
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

    cloud := config.GetUserConfig("cloud")
    if err := cloud.UploadFile(bucket.BucketName, file); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to upload file": err.Error()})
    } else {
        c.String(http.StatusCreated, "")
    }
}

func downloadFile(c *gin.Context) {
    var f File

    if err := c.ShouldBindUri(&f); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }

    cloud := config.GetUserConfig("cloud")
    if object,err := cloud.DownloadFile(f.BucketName, f.FileName); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to download file": err.Error()})
    } else {
        if stat,err := object.Stat(); err != nil {
            c.JSON(http.StatusNotFound, gin.H{"failed to get file stat": err.Error()})
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
        c.String(http.StatusBadRequest, err.Error())
        return
    }

    cloud := config.GetUserConfig("cloud")
    if err := cloud.RemoveFile(f.BucketName, f.FileName); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to delete file": err.Error()})
    } else {
        c.String(http.StatusOK, "")
    }
}

func getUserKey(c *gin.Context) {
    var bucket Bucket
    if err := c.ShouldBindUri(&bucket); err != nil {
        c.String(http.StatusBadRequest, err.Error())
        return
    }
    admin := config.GetUserConfig("admin")
    if key,err := admin.GetUserOfBucket(bucket.BucketName); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"failed to get user key": err.Error()})
    } else {
        config.SetBucketConfig(c.Query("um"), "config.yaml", map[string]string{bucket.BucketName: "rw"})
        c.JSON(http.StatusOK, key)
    }
}

func setUserKey(c *gin.Context) {
    ssl,_ := strconv.ParseBool(c.DefaultPostForm("ssl", "false"))
    userInfo := backends.S3UserInfo {
        Endpoint: c.PostForm("endpoint"),
        AccessKey: c.PostForm("accesskey"),
        SecretKey: c.PostForm("secretkey"),
        SSL: ssl,
        Location: c.DefaultPostForm("location", "us-east-1"),
    }
    config.Users[c.Param("user")] = &userInfo
    utils.SaveConfig(yamlConfigFile, config)
    c.JSON(http.StatusOK, "")
}
