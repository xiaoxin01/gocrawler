package service

import (
	"testing"

	"gocrawler/model"

	"github.com/stretchr/testify/assert"
)

func TestAlert(t *testing.T) {
	title := "标题"
	description := "描述"
	model.InitConfig("../")
	model.InitField("../")

	success := Alert(title, description)

	assert.True(t, success)
}
