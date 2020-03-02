package service

import "fmt"

var (
	alertService ialert
)

// InitAlert initial alert type and key
func InitAlert(alertType string, alertKey string) {
	switch alertType {
	case "wechat":
		useWechatAlert(alertKey)
	default:
		useScAlert(alertKey)
	}
}

// Alert send alert
func Alert(title string, content string) bool {
	return alertService.Alert(title, content)
}

func useScAlert(key string) {
	alertService = scAlert{Key: key, URL: fmt.Sprintf("https://sc.ftqq.com/%s.send", key)}
}

func useWechatAlert(key string) {
	alertService = wechatAlert{Key: key, URL: fmt.Sprintf("http://wechat.supperxin.com/message/%s", key)}
}
