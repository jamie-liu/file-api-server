package backends

import (
    "fmt"
    "github.com/golang/glog"
    "github.com/minio/minio-go"
    "mime/multipart"
    "path/filepath"
)

func initConnection(endpoint, accessKey, secretKey string, secure bool) (*minio.Client, error) {
    conn, err := minio.New(endpoint, accessKey, secretKey, secure)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func initConnectionWithRegion(endpoint, accessKey, secretKey string, secure bool, location string) (*minio.Client, error) {
    conn, err := minio.NewWithRegion(endpoint, accessKey, secretKey, secure, location)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func (user *S3UserInfo) CreateBucket(bucketName string) error {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
    if err != nil {
        glog.Errorln(err)
        return err
    }

    if err := conn.MakeBucket(bucketName, user.Location); err != nil {
        glog.Errorln(err)
        return err
    }
    return nil
}

func (user *S3UserInfo) ListBuckets() ([]minio.BucketInfo, error) {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
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

func (user *S3UserInfo) RemoveBucket(bucketName string) error {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
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

func (user *S3UserInfo) ListBucket(bucketName, prefix string, recursive bool) ([]FileStat, error) {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
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

        objectCh := conn.ListObjects(bucketName, prefix, recursive, doneCh)
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

func (user *S3UserInfo) GetBucketPolicy(bucketName string) (string, error) {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
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

func (user *S3UserInfo) UploadLocalFile(bucketName, objectName, filePath string) error {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
    if err != nil {
        glog.Errorln(err)
        return err
    }

    if exists, err := conn.BucketExists(bucketName); err != nil {
        glog.Errorln(err)
        return err
    } else if !exists {
        err = conn.MakeBucket(bucketName, user.Location)
        if err != nil {
            glog.Errorln(err)
            return err
        }
    }

    if n, err := conn.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{}); err != nil {
        glog.Errorln(err)
        return err
    } else {
        glog.Infof("Successfully uploaded %s of size %d\n", objectName, n)
        return nil
    }
}

func (user *S3UserInfo) UploadFile(bucketName string, file *multipart.FileHeader) error {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
    if err != nil {
        glog.Errorln(err)
        return err
    }

    if exists, err := conn.BucketExists(bucketName); err != nil {
        glog.Errorln(err)
        return err
    } else if !exists {
        err = conn.MakeBucket(bucketName, user.Location)
        if err != nil {
            glog.Errorln(err)
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
        glog.Errorln(err)
        return err
    } else {
        glog.Infof("Successfully uploaded %s of size %d\n", filename, n)
        return nil
    }
}

func (user *S3UserInfo) DownloadFile(bucketName, objectName string) (*minio.Object, error) {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
    if err != nil {
        glog.Errorln(err)
        return nil, err
    }

    if exists, err := conn.BucketExists(bucketName); err != nil {
        glog.Errorln(err)
        return nil, err
    } else if !exists {
        return nil, fmt.Errorf("Bucket %s doesn't exist\n", bucketName)
    }

    //if err := conn.FGetObject(bucketName, objectName, filePath, minio.GetObjectOptions{}); err != nil {
    return conn.GetObject(bucketName, objectName, minio.GetObjectOptions{})
}

func (user *S3UserInfo) RemoveFile(bucketName, objectName string) error {
    conn, err := initConnectionWithRegion(user.Endpoint, user.AccessKey, user.SecretKey, user.SSL, user.Location)
    if err != nil {
        glog.Errorln(err)
        return err
    }

    if exists, err := conn.BucketExists(bucketName); err != nil {
        return err
    } else if !exists {
        return fmt.Errorf("Bucket %s doesn't exist\n", bucketName)
    }

    if err := conn.RemoveObject(bucketName, objectName); err != nil {
        glog.Errorln(err)
        return err
    }
    return nil
}
