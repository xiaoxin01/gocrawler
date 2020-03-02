package service

import (
	"reflect"
	"testing"

	"gocrawler/model"

	"github.com/stretchr/testify/assert"
)

func TestUseService(t *testing.T) {
	t.Run("test use sc alert", func(t *testing.T) {
		useScAlert("")

		assert.Equal(t, reflect.TypeOf(scAlert{}), reflect.TypeOf(alertService))
	})

	t.Run("test use wechat alert", func(t *testing.T) {
		useWechatAlert("")

		assert.Equal(t, reflect.TypeOf(wechatAlert{}), reflect.TypeOf(alertService))
	})
}

func TestSendAlert(t *testing.T) {
	t.Run("test send wechat alert", func(t *testing.T) {
		title := "标题"
		description := "描述"
		model.InitConfig("../")
		model.InitField("../")
		InitAlert(model.AlertType, model.AlertKey)

		success := alertService.Alert(title, description)

		assert.True(t, success)
	})
}
