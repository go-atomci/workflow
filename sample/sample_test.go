package sample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSamplePipeline(t *testing.T) {
	err := SamplePipeline()
	assert.Nil(t, err)
}
