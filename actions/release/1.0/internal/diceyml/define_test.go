// Copyright (c) 2021 Terminus, Inc.
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

package diceyml

import (
	"testing"
)
func TestValidImageName(t *testing.T) {
	validRepoNames := []string{
		"docker/docker:0.0.1",
		"index.docker.io/debian:latest",
		"127.0.0.1:5000/debian:v1.0.1",
		"thisisthesongthatneverendsitgoesonandonandonthisisthesongthatnev:test-fix",
		"docker.io/1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
	}
	invalidRepoNames := []string{
		"DOCKER/docker",
		"https://github.com/docker/docker",
		"docker/Docker",
		"-docker/docker",
		"docker///docker",
		"docker.io/docker/Docker",
		"1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
		"docker/docker:v0.0.1 ",
		"docker/docker:v0.0.1 \u200b",
	}

	for _, name := range invalidRepoNames {
		err := ValidImageName(name)
		if err == nil {
			t.Fatalf("Expected invalid repo name for %q", name)
		}
	}

	for _, name := range validRepoNames {
		err := ValidImageName(name)
		if err != nil {
			t.Fatalf("Error parsing repo name %s, got: %q", name, err)
		}
	}
}