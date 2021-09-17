// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package oss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClient_DeleteRemote(t *testing.T) {
	endpoint := os.Getenv("TEST_OSS_ENDPOINT")
	accessKeyID := os.Getenv("TEST_OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("TEST_OSS_ACCESS_KEY_SECRET")
	bucket := os.Getenv("TEST_OSS_BUCKET")
	remote := os.Getenv("TEST_OSS_REMOTE_PREFIX")
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucket == "" || remote == "" {
		return
	}

	client, err := New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		t.Fatalf("failed to oss.New: %v", err)
	}

	// delete empty dir
	obj := defaultObject{
		bucket: bucket,
		remote: filepath.Join(remote, "del-empty"),
		local:  "",
	}
	if err = client.DeleteRemoteRecursively(obj); err != nil {
		t.Fatalf("failed to delete empty dir, obj: %+v: %v", obj, err)
	}

	// delete non-empty
	obj.remote = filepath.Join(remote, "del-non-empty")
	if err = client.DeleteRemoteRecursively(obj); err != nil {
		t.Fatalf("faield to delete non-empty dir recursively, obj: %+v: %v", obj, err)
	}
}
