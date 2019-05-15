package utils

import (
    "io/ioutil"
    "fmt"
    "net/http"
    "encoding/json"
)

type Response struct {
    Err  int    `json:"err"`
    Msg  string `json:"msg"`
    User string `json:"user"`
}

func VerifyToken(url, token string) (*Response, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s?tk=%s", url, token), nil)
    req.Header.Set("Content-Type", "application/json;charset=UTF-8")
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        return nil, err
    } else {
        var ret Response
        if err := json.Unmarshal(body, &ret); err != nil {
            return nil, err
        }
        return &ret, nil
    }
}
