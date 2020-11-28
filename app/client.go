package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	Endpoint string = "https://notify-bot.line.me"
)

func Call(method string, url string, data url.Values, token string) ([]byte, error) {
	client := &http.Client{}
	res := &http.Request{}

	if data == nil {
		res, _ = http.NewRequest(method, url, nil)
	} else {
		res, _ = http.NewRequest(method, url, strings.NewReader(data.Encode()))
	}

	res.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if len(token) != 0 {
		res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := client.Do(res)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode > http.StatusAccepted {
		return body, err
	}

	return body, nil
}
