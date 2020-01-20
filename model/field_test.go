package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	t.Run("test get value from operator: Func", func(t *testing.T) {
		field := Field{Operator: "Func", Parameter: "time.Now().Unix()"}

		time := time.Now().Unix()
		value, ok := field.GetValue(nil)

		assert.True(t, ok)
		assert.LessOrEqual(t, time, value)
	})
}
