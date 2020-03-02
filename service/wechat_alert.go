package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type wechatAlert struct {
	Key string
	URL string
}

// Alert send alert
func (alert wechatAlert) Alert(title string, content string) bool {
	success := true
	alertURL := fmt.Sprintf("%s?title=%s&content=%s", alert.URL, url.QueryEscape(title), url.QueryEscape(content))
	resp, err := http.Get(alertURL)
	if err != nil {
		fmt.Println(err)
		success = false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		success = false
	}

	fmt.Println(string(body))

	return success
}
