package utils

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "github.com/golang/glog"
)

func LoadConfig(path string, conf interface{}) error {
    if data, err := ioutil.ReadFile(path); err != nil {
        return err
    } else {
        glog.Infof("---dump yaml file---:\n%s\n", string(data))
        return yaml.Unmarshal(data, conf)
    }
}

func SaveConfig(path string, conf interface{}) error {
    if data, err := yaml.Marshal(conf); err != nil {
        return err
    } else {
        glog.Infof("---dump yaml file---:\n%s\n", string(data))
        return ioutil.WriteFile(path, data, 0644)
    }
}
