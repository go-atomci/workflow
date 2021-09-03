package sample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSamplePipeline(t *testing.T) {
	addr := "http://10.10.1.150:8091"
	user := "admin"
	token := "11e0cd1fd42d54fc8e47e9224144179b7"
	err := SamplePipeline(addr, user, token)
	assert.Nil(t, err)
}
