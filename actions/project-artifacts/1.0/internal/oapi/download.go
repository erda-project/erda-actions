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

package oapi

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/project-artifacts/1.0/internal/config"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

// DownloadArtifact downloads the project artifact with the given releaseID
func DownloadArtifact(cfg *config.Config, releaseID string) error {
	// Create the output file
	out, err := os.Create("artifact.zip")
	if err != nil {
		return errors.Wrap(err, "failed to create artifact.zip")
	}
	defer out.Close()

	// Perform the download using httpclient
	response, err := httpclient.New().
		Get(cfg.OapiHost).
		Path(fmt.Sprintf("/api/releases/%s/actions/download", releaseID)).
		Header("Authorization", cfg.OapiToken).
		Do().
		Body(out)

	if err != nil {
		return errors.Wrap(err, "failed to download artifact")
	}

	if !response.IsOK() {
		return errors.Errorf("failed to download artifact, status code: %d", response.StatusCode())
	}

	return nil
}