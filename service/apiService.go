package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func GetJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 5 * time.Second}

	response, err := myClient.Get(url)
	if err != nil {
		return err
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}

	return nil
}
