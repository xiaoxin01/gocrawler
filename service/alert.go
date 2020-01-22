package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gocrawler/model"
)

// Alert send alert to your device
func Alert(title string, description string) bool {
	success := true
	url := fmt.Sprintf("https://sc.ftqq.com/%s.send", model.AlertKey)
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("text=%s&desp=%s", title, description)))
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
