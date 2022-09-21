package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type response struct {
	Status bool
}

func TestParseJsonFromFilePath(t *testing.T) {
	var s response
	err := JSONModelFromFilePath("test.json", &s)

	assert.Nil(t, err)
	assert.True(t, s.Status)
}

func TestJSONStringFromFilePath(t *testing.T) {
	data, err := JSONStringFromFilePath("test.json")
	assert.Nil(t, err)
	assert.Equal(t, `{
  "status": true
}`, data)
}
