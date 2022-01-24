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

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/erda-create-custom-addon/1.0/internal/config"
	"github.com/erda-project/erda-actions/actions/erda-create-custom-addon/1.0/internal/oapi"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func main() {
	logrus.Infoln("Custom Addon begins to work")
	addon, err := oapi.New(config.Get()).Create()
	if err != nil {
		_ = metawriter.WriteSuccess(false)
		logrus.WithError(err).Fatalln("failed to Create custom addon")
	}
	logrus.
		WithFields(map[string]interface{}{
			"name":              addon.Name,
			"tag":               addon.Tag,
			"configs":           string(addon.Config),
			"instanceID":        addon.InstanceId,
			"routingInstanceID": addon.RealInstanceId,
		}).
		Infoln("the addon info")
	_ = metawriter.WriteKV("name", addon.Name)
	_ = metawriter.WriteKV("tag", addon.Tag)
	_ = metawriter.WriteKV("configs", string(addon.Config))
	_ = metawriter.WriteLink("addonInstanceID", addon.InstanceId)
	_ = metawriter.WriteKV("routingInstanceID", addon.RealInstanceId)
	_ = metawriter.WriteSuccess(true)
}
