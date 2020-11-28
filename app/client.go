package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Call(method string, url string, data url.Values, token string) ([]byte, error) {
	client := &http.Client{}

	res, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	res.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if token != "" {
		res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := client.Do(res)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
