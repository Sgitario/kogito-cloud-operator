// Copyright 2019 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FromImageTagToImage(t *testing.T) {
	buildImage := "quay.io/vajain/kogito-cloud-operator:1.0"
	image := FromImageTagToImage(buildImage)
	assert.NotNil(t, image)
	assert.Equal(t, "quay.io", image.Domain)
	assert.Equal(t, "vajain", image.Namespace)
	assert.Equal(t, "kogito-cloud-operator", image.Name)
	assert.Equal(t, "1.0", image.Tag)

}
