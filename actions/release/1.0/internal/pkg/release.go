package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

const (
	releaseBranchPrefix = "release/"
)

func genReleaseRequest(cfg *conf.Conf) *apistructs.ReleaseCreateRequest {
	labels := make(map[string]string, len(cfg.Labels)+5)
	// insert user defined label
	for k, v := range cfg.Labels {
		labels[k] = v
	}
	labels["gitRepo"] = cfg.GittarRepo
	labels["gitBranch"] = cfg.GittarBranch
	labels["gitCommitId"] = cfg.GittarCommitID
	labels["gitCommitMessage"] = cfg.GittarMessage

	release := &apistructs.ReleaseCreateRequest{
		ReleaseName:     cfg.GittarBranch,
		Labels:          labels,
		OrgID:           cfg.OrgID,
		ProjectID:       cfg.ProjectID,
		ApplicationID:   cfg.AppID,
		UserID:          cfg.DiceOperatorID,
		ClusterName:     cfg.ClusterName,
		ProjectName:     cfg.ProjectName,
		ApplicationName: cfg.AppName,
		CrossCluster:    cfg.CrossCluster,
	}

	//  如果是TAG，则插入到version里
	if cfg.ReleaseTag != "" {
		release.Version = cfg.ReleaseTag
	}

	return release
}

// push release info to dicehub
func pushRelease(cfg conf.Conf, req *apistructs.ReleaseCreateRequest) (string, error) {
	var releaseCreateResp apistructs.ReleaseCreateResponse
	err := retry.DoWithInterval(func() error {
		request := httpclient.New(httpclient.WithCompleteRedirect()).Post(cfg.DiceOpenapiPrefix).
			Path("/api/releases").
			Header("Authorization", cfg.CiOpenapiToken)
		if strings.Compare(cfg.DiceVersion, "3.12") >= 0 || cfg.Base64Switch {
			reqBytes, err := json.Marshal(req)
			if err != nil {
				return err
			}
			request = request.Header("base64-encoded-request-body", "true").
				RawBody(bytes.NewBufferString(base64.StdEncoding.EncodeToString(reqBytes)))
		} else {
			request = request.JSONBody(req)
		}
		resp, err := request.Do().JSON(&releaseCreateResp)
		if err != nil {
			return err
		}
		if !resp.IsOK() {
			return errors.Errorf("failed to push release, status code: %d, response body: %v", resp.StatusCode(), string(resp.Body()))
		}
		if !releaseCreateResp.Success {
			return errors.Errorf(releaseCreateResp.Error.Msg)
		}
		return nil
	}, 2, time.Second*1)

	if err != nil {
		return "", err
	}

	return releaseCreateResp.Data.ReleaseID, nil
}
