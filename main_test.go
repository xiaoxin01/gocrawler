package main

import "testing"

import "regexp"

import "github.com/magiconair/properties/assert"

func TestRegex(t *testing.T) {
	t.Run("test regex match", func(t *testing.T) {
		url := "details.php?id=123456&ref="
		finds := regexp.MustCompile("id=(\\d+)").FindStringSubmatch(url)
		assert.Equal(t, len(finds), 2)
		assert.Equal(t, finds[1], "123456")
	})
}
