/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
