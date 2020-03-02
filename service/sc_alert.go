package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type scAlert struct {
	Key string
	URL string
}

// Alert send alert
func (alert scAlert) Alert(title string, content string) bool {
	success := true
	resp, err := http.Post(alert.URL,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("text=%s&desp=%s", title, content)))
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
