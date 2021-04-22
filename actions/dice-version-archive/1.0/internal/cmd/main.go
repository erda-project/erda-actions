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

package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-version-archive/1.0/internal/archive"
	"github.com/erda-project/erda-actions/actions/dice-version-archive/1.0/internal/config"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func main() {
	logrus.Infoln("Dice Version Archive start working")

	// read VERSION file
	versionFile := filepath.Join(config.Workdir(), "VERSION")
	version := new(archive.Version)
	if err := version.Read(versionFile); err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to read version file: %v", err)
	}

	api := archive.NewAccessAPI(
		config.OpenapiPrefix(),
		config.OpenapiToken(),
		strconv.FormatUint(config.OrgID(), 10),
		config.ProjectName(),
		config.DstApplicationName(),
		config.ReleaseID(),
	)

	// read dice.yml
	diceyaml := new(archive.DiceYaml)
	if err := diceyaml.ReadFromDiceHub(api); err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to read dice.yml: %v", err)
	}
	deployableContent, err := diceyaml.Deployable()
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to read dice: %v", err)
	}

	// read migrations scripts files
	logrus.Infoln("read migration scripts")
	scripts, err := archive.ReadScripts(config.Workdir(), config.MigrationsPathFromSrcRepoRoot())
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Warn: "failed to read migrations scripts, no script is archived"})
		logrus.Warnf("failed to read migration scripts: %v", err)
	}

	// create new branch, commit, merge request in src repo by gittar handler
	logrus.Infoln("do archiving: create new branch, commit, merge request")
	gittar := archive.NewGittar(api)

	createBranchPayload := archive.CreateBranchPayload{
		Name: config.DstRepoBranch(),
		Ref:  config.DstRepoRefBranch(),
	}
	logrus.Infof("create branch %s refered from master on dst repo", createBranchPayload.Name)
	if err = gittar.CreateBranch(&createBranchPayload); err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to CreateBranch on DstReop: %v", err)
	}

	logrus.Infoln("create commit on dst repo")
	message := fmt.Sprintf("archive dice.yml and migrations from %s/%s repo", config.ProjectName(), config.ApplicationName())
	createCommitPayload := archive.CreateCommitPayload{
		Message: message,
		Branch:  createBranchPayload.Name,
		Actions: []*archive.CreateCommitPayloadAction{{
			Action:   archive.ActionAdd,
			Content:  deployableContent,
			Path:     filepath.Join(version.String(), config.DiceYmlPathFromDstRepoVersionDir()),
			PathType: archive.PathTypeBlob,
		}, {
			Action:   archive.ActionAdd,
			Content:  version.String(),
			Path:     "version",
			PathType: archive.PathTypeBlob,
		}},
	}
	for _, script := range scripts {
		action := archive.CreateCommitPayloadAction{
			Action:   archive.ActionAdd,
			Content:  string(script.Content),
			Path:     filepath.Join(version.String(), config.MigrationPathFromDstRepoVersionDir, script.NameFromService),
			PathType: archive.PathTypeBlob,
		}
		createCommitPayload.Actions = append(createCommitPayload.Actions, &action)
	}
	if err = gittar.CreateCommit(&createCommitPayload); err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to CreateCommit on DesRepo: %v", err)
	}

	logrus.Infoln("create merge request on dst repo")
	createMergeRequestPayload := archive.CreateMergeRequestPayload{
		Title: message,
		Description: fmt.Sprintf("this merge request is created from [pipeline-%s](/workBench/projects/%v/apps/%v/pipeline/%s)",
			config.PipelineID(), config.ProjectID(), config.ApplicationID(), config.PipelineID()),
		AssigneeID:         config.MRProcessor(),
		SourceBranch:       createBranchPayload.Name,
		TargetBranch:       createBranchPayload.Ref,
		RemoveSourceBranch: true,
	}
	id, err := gittar.CreateMergeRequest(&createMergeRequestPayload)
	if err != nil {
		_ = metawriter.Write(map[string]interface{}{config.Success: false, config.Err: err})
		logrus.Fatalf("failed to CreateMergeRequest on DesRepo: %v", err)
	}

	_ = metawriter.Write(map[string]interface{}{config.Success: true, config.MrID: id})
	logrus.Infoln("archive complete.")
}
