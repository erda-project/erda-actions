package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type Conf struct {
	ApplicationName string `env:"ACTION_APPLICATION_NAME" required:"true"`
	Branch          string `env:"ACTION_BRANCH" required:"true"`
	CheckCommit     string `env:"ACTION_CHECK_COMMIT" required:"false"`

	// sys
	DiceProjectID    string `env:"DICE_PROJECT_ID" required:"true"`
	DiceOpenapiAddr  string `env:"DICE_OPENAPI_ADDR" required:"true"`
	DiceOpenapiToken string `env:"DICE_OPENAPI_TOKEN" required:"true"`
}

var conf Conf

func main() {
	envconf.MustLoad(&conf)

	// http client
	hc := httpclient.New(
		httpclient.WithCompleteRedirect(),
		httpclient.WithTimeout(time.Second, time.Second*3),
	)

	appID, err := getAppID(hc, conf.ApplicationName)
	if err != nil {
		panic(err)
	}

	// get release
	release, err := getRelease(hc, appID)
	if err != nil {
		echoMeta("Error", err.Error())
		panic(err)
	}

	echoMeta("release_id", release.ReleaseID)
	echoMeta("release_name", release.Version)
	echoMeta("release_branch", release.Labels["gitBranch"])
	echoMeta("release_commit", release.Labels["gitCommitId"])
	echoMeta("release_commit_message", release.Labels["gitCommitMessage"])

	// check commit
	if conf.CheckCommit != "" {
		if conf.CheckCommit != release.Labels["gitCommitId"] {
			echoMeta("Error", "commit id not match")
			panic("commit id not match")
		}
	}
}

func echoMeta(k, v string) {
	fmt.Printf("action meta: %s=%s\n", k, v)
}

func getAppID(hc *httpclient.HTTPClient, name string) (string, error) {
	var resp apistructs.ApplicationListResponse
	r, err := hc.Get(conf.DiceOpenapiAddr).Path("/api/applications").
		Param("projectId", conf.DiceProjectID).
		Param("name", name).
		Param("pageNo", "1").
		Param("pageSize", "1").
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)
	if err != nil {
		return "", err
	}
	if !r.IsOK() || !resp.Success {
		return "", fmt.Errorf(resp.Error.Msg)
	}
	if resp.Data.Total == 0 || len(resp.Data.List) == 0 {
		return "", fmt.Errorf("application not found")
	}
	return strconv.FormatUint(resp.Data.List[0].ID, 10), nil
}

func getRelease(hc *httpclient.HTTPClient, appID string) (*apistructs.ReleaseData, error) {
	var resp apistructs.ReleaseListResponse
	// fetch release
	r, err := hc.Get(conf.DiceOpenapiAddr).
		Path("/api/releases").
		Param("projectId", conf.DiceProjectID).
		Param("applicationId", appID).
		Param("branchName", conf.Branch).
		Param("pageNo", "1").
		Param("pageSize", "1").
		Header("Authorization", conf.DiceOpenapiToken).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() || !resp.Success {
		return nil, fmt.Errorf(resp.Error.Msg)
	}
	if resp.Data.Total == 0 || len(resp.Data.Releases) == 0 {
		return nil, fmt.Errorf("release not found")
	}
	return &resp.Data.Releases[0], nil
}
