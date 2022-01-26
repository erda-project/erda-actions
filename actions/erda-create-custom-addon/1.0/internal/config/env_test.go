// Copyright (c) 2022 Terminus, Inc.
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

package config_test

import (
	"os"
	"testing"

	"github.com/erda-project/erda-actions/actions/erda-create-custom-addon/1.0/internal/config"
)

func TestFindEnvLiteral(t *testing.T) {
	var (
		evns = []string{
			"${my_name}",
			"${gateway}/url",
			"my${parent}child",
			"${}",
		}
	)
	for _, s := range evns {
		env, indexStart, indexEnd, err := config.FindEnvLiteral(s)
		if err != nil {
			t.Fatal(err, s)
		}
		t.Log(env, indexStart, indexEnd, s[indexStart:indexEnd])
	}
	env, indexStart, indexEnd, err := config.FindEnvLiteral(`some${info
}in the multilines`)
	if err == nil {
		t.Fatal("err not be nil")
	}
	t.Log(env, indexStart, indexEnd, err)
}

func TestInterpolate(t *testing.T) {
	var m = map[string]string{
		"GATEWAY_URL":           "test-gateway.app.terminus.io",
		"DICE_PROJECT_NAME":     "trantor-auto-deploy",
		"DICE_APPLICATION_NAME": "trantor-metastore",
		"META_STORE_URL":        "${GATEWAY_URL}/${DICE_PROJECT_NAME}/${DICE_APPLICATION_NAME}/${SERVICE_NAME}",
	}
	if err := os.Setenv("SERVICE_NAME", "metastore-runtime"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := config.Interpolate(m); err != nil {
		t.Fatalf("failed to Interpolate: %v", err)
	}
	t.Logf("%+v", m)
	if m["META_STORE_URL"] != "test-gateway.app.terminus.io/trantor-auto-deploy/trantor-metastore/metastore-runtime" {
		t.Fatal("failed to interpolate")
	}
}
