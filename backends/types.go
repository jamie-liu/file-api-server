package backends

import (
    "time"
    "github.com/file-api-server/utils"
    "github.com/golang/glog"
)

type S3UserInfo struct {
    Endpoint        string      `yaml: "endpoint" json:"endpoint"`
    AccessKey       string      `yaml: "accesskey" json:"accesskey"`
    SecretKey       string      `yaml: "secretkey" json:"secretkey"`
    SSL             bool        `yaml: "ssl" json:"ssl"`
    Location        string      `yaml: "location" json:"location"`
}

type Config struct {
    URL             string                          `yaml:"cm,omitempty"`
    Users           map[string]*S3UserInfo          `yaml:"users,omitempty"`
    Buckets         map[string]map[string]string    `yaml:"buckets,omitempty"`
}

type FileStat struct {
    Name            string      `json:"name,omitempty"`
    Size            int64       `json:"size,omitempty"`
    LastModified    time.Time   `json:"lastModified,omitempty"`
}

type S3Key struct {
    AccessKey       string      `json:"access,omitempty"`
    SecretKey       string      `json:"secret,omitempty"`
    Endpoint        string      `json:"endpoint,omitempty"`
}

func (conf *Config) GetUserConfig(userName string) *S3UserInfo {
    if _,ok := conf.Users[userName]; !ok {
        glog.Fatalf("User %s config doesn't exsit in yaml", userName)
    }
    return conf.Users[userName]
}

func (conf *Config) SetUserConfig(userName, privateFile, publicFile, configFile string, user *S3UserInfo) {
    conf.Users[userName] = user
    for _,v := range conf.Users {
        v.AccessKey,_ = utils.RsaEncrypt(v.AccessKey, publicFile)
        v.SecretKey,_ = utils.RsaEncrypt(v.SecretKey, publicFile)
    }
    utils.SaveConfig(configFile, conf)
    for _,v := range conf.Users {
        v.AccessKey,_ = utils.RsaDecrypt(v.AccessKey, privateFile)
        v.SecretKey,_ = utils.RsaDecrypt(v.SecretKey, privateFile)
    }
}

func (conf *Config) SetBucketConfig(um, configFile string, bucket map[string]string) {
    if conf.Buckets == nil {
        conf.Buckets = make(map[string]map[string]string)
    }
    conf.Buckets[um] = bucket
    utils.SaveConfig(configFile, conf)
}

func (conf *Config) GetBucketConfig(um, bucket string) (string, bool) {
    if _,ok := conf.Buckets[um]; ok {
        return conf.Buckets[um][bucket], ok
    }
    return "", false
}
