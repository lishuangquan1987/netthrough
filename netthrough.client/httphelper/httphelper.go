package httphelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string) (string, error) {
	rsp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fail to request:%s,reason:%v", url, err)
		return "", err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("read body error,url:%s,reason:%v", url, err)
		return "", err
	}
	return string(body), nil
}

func GetObj(url string, obj interface{}) error {
	rsp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fail to request:%s,reason:%v", url, err)
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("read body error,url:%s,reason:%v", url, err)
		return err
	}

	return json.Unmarshal(body, obj)
}

func Post(url string, data interface{}) (string, error) {
	var requstBytes []byte = nil
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("convert to json fail.url:%s", url)
			return "", err
		}
		requstBytes = bytes
	}
	var rsp *http.Response
	var err error
	if requstBytes != nil {
		rsp, err = http.Post(url, "application/json", bytes.NewBuffer(requstBytes))
	} else {
		rsp, err = http.Post(url, "application/json", nil)
	}

	if err != nil {
		fmt.Printf("fail to request url:%s,reason:%v", url, err)
		return "", err
	}
	defer rsp.Body.Close()
	responseBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("read body error,url%s,reason:%v", url, err)
		return "", err
	}
	return string(responseBytes), nil
}
func PostObj(url string, data interface{}, obj interface{}) error {
	var requstBytes []byte = nil
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("convert to json fail.url:%s", url)
			return err
		}
		requstBytes = bytes
	}
	var rsp *http.Response
	var err error
	if requstBytes != nil {
		rsp, err = http.Post(url, "application/json", bytes.NewBuffer(requstBytes))
	} else {
		rsp, err = http.Post(url, "application/json", nil)
	}

	if err != nil {
		fmt.Printf("fail to request url:%s,reason:%v", url, err)
		return err
	}
	defer rsp.Body.Close()
	responseBytes, err := ioutil.ReadAll(rsp.Body)
	fmt.Printf("PostObj Receive bytes:%d", len(responseBytes))
	if err != nil {
		fmt.Printf("read body error,url%s,reason:%v", url, err)
		return err
	}
	return json.Unmarshal(responseBytes, obj)
}
