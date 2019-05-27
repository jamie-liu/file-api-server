package backends

import (
    "time"
    "github.com/file-api-server/utils"
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
