package backends

import (
    "time"
)

type S3UserInfo struct {
    Endpoint        string      `yaml: "endpoint"`
    AccessKey       string      `yaml: "accesskey"`
    SecretKey       string      `yaml: "secretkey"`
    SSL             bool        `yaml: "ssl"`
    Location        string      `yaml: "location"`
}

type Config struct {
    URL             string                     `yaml:"cm"`
    Users           map[string]*S3UserInfo     `yaml:"users"`
    Buckets         map[string][]string        `yaml:"buckets"`
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
