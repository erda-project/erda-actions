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

package repo

import (
	"io/ioutil"
	"testing"

	"github.com/erda-project/erda/pkg/parser/diceyml"
	"sigs.k8s.io/yaml"
)

func TestPatch(t *testing.T) {
	data, err := ioutil.ReadFile("test_dice.yml")
	if err != nil {
		t.Fatal(err)
	}
	deployable, err := diceyml.NewDeployable(data, diceyml.WS_PROD, false)
	if err != nil {
		t.Fatal(err)
	}
	obj := deployable.Obj()
	patchSecurityContextPrivileged(obj, "cluster-agent")
	clusterAgent := obj.Services["cluster-agent"]
	if clusterAgent == nil {
		t.Fatal("cluster-agent can not be nil")
	}
	if clusterAgent.K8SSnippet == nil {
		t.Fatal("cluster-agent.K8SSnippet can not be nil")
	}
	if clusterAgent.K8SSnippet.Container == nil {
		t.Fatal("cluster-agent.K8SSnippet.Container can not be nil")
	}
	if clusterAgent.K8SSnippet.Container.SecurityContext == nil {
		t.Fatal("cluster-agent.K8SSnippet.Container.SecurityContext can not be nil")
	}
	if privileged := clusterAgent.K8SSnippet.Container.SecurityContext.Privileged; privileged == nil || !*privileged {
		t.Fatal("cluster-agent.K8SSnippet.Container.SecurityContext.Privileged can not be nil or false")
	}

	data, err = yaml.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
