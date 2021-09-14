package appCreate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/app-create/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

const CloneAddr = "_repo"

func handleAPIs() error {
	err := os.Mkdir(CloneAddr, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir %s error: %v", CloneAddr, err)
	}

	// check application type
	switch conf.ApplicationType() {
	case string(apistructs.ApplicationModeService), string(apistructs.ApplicationModeLibrary),
		string(apistructs.ApplicationModeBigdata), string(apistructs.ApplicationModeAbility),
		string(apistructs.ApplicationModeMobile), string(apistructs.ApplicationModeApi):
	default:
		return errors.New("invalid request, mode is invalid, just support LIBRARY, SERVICE, MOBILE")
	}

	// if not external repo, git clone code
	repo := conf.ApplicationGitRepo()
	if !conf.IsExternalRepo() {
		if len(conf.ApplicationGitUsername()) > 0 && len(conf.ApplicationGitPassword()) > 0 {
			var splitValue []string
			var preFix = ""
			if strings.HasPrefix(repo, "http://") {
				splitValue = strings.SplitN(repo, "http://", 2)
				preFix = "http://"
			} else if strings.HasPrefix(repo, "https://") {
				splitValue = strings.SplitN(repo, "https://", 2)
				preFix = "https://"
			} else {
				return fmt.Errorf("application git repo addr just support http")
			}

			if len(splitValue) != 2 {
				return fmt.Errorf("application git repo addr append token error")
			}
			repo = preFix + fmt.Sprintf("%s:%s@", url.QueryEscape(conf.ApplicationGitUsername()), url.QueryEscape(conf.ApplicationGitPassword())) + splitValue[1]
		}

		logrus.Infof("start git clone url %s", conf.ApplicationGitRepo())
		err = runCommand(fmt.Sprintf("git clone %s %s", repo, CloneAddr), true)
		if err != nil {
			return fmt.Errorf("run git clone error: %v", err)
		}
		logrus.Infof("end git clone")

		err = os.Chdir(CloneAddr)
		if err != nil {
			return fmt.Errorf("chdir %s error: %v", CloneAddr, err)
		}

		_ = runCommand(fmt.Sprintf("for b in `git branch -r | grep -v -- '->'`; do git branch --track ${b##origin/} $b; done\n"), true)
	}

	// check application exit
	logrus.Infof("start get application %s", conf.ParamsApplicationName())
	existApp, err := checkAppExist()
	if err != nil {
		return err
	}
	logrus.Infof("end get application")
	// application exist return appID and info
	if existApp != nil {
		return storeMetaFile(strconv.FormatUint(existApp.ID, 10), true)
	}

	logrus.Infof("start to create application %s", conf.ParamsApplicationName())
	// create application
	var req apistructs.ApplicationCreateRequest
	req.Name = conf.ParamsApplicationName()
	req.Mode = apistructs.ApplicationMode(conf.ApplicationType())
	req.ProjectID = conf.ProjectId()
	req.IsExternalRepo = conf.IsExternalRepo()
	if req.IsExternalRepo {
		req.RepoConfig = &apistructs.GitRepoConfig{
			Password: conf.ApplicationGitPassword(),
			Username: conf.ApplicationGitUsername(),
			Url:      conf.ApplicationGitRepo(),
			Type:     "general",
		}
	}
	dbApplication, err := createApplication(req)
	if err != nil {
		return err
	}
	logrus.Infof("end to create application %s", conf.ParamsApplicationName())

	dbApplication, err = getApplication(dbApplication.ID, dbApplication.Creator)
	if err != nil {
		return err
	}

	// if not external repo push code to repo
	if !conf.IsExternalRepo() {
		logrus.Infof("start push code to application %s", conf.ParamsApplicationName())
		err := runCommand(fmt.Sprintf("git remote add app_create_dice https://%s:%s@%s", url.QueryEscape(conf.GittarUsername()), url.QueryEscape(conf.GittarPassword()), dbApplication.GitRepoNew), true)
		if err != nil {
			return fmt.Errorf("git remote add error: %v", err)
		}

		if len(conf.ApplicationGitBranchs()) <= 0 {
			err = runCommand(fmt.Sprintf("git push -u app_create_dice --all"), true)
			if err != nil {
				return fmt.Errorf("git remote add error: %v", err)
			}
		} else {
			for _, v := range conf.ApplicationGitBranchs() {
				err = runCommand(fmt.Sprintf("git push -u app_create_dice %v", v), true)
				if err != nil {
					return fmt.Errorf("git remote add error: %v", err)
				}
			}
		}

		logrus.Infof("end push code to application")
	}

	return storeMetaFile(strconv.FormatUint(dbApplication.ID, 10), false)
}

func runCommand(cmd string, hideStdout bool) error {
	command := exec.Command("/bin/bash", "-c", cmd)
	if hideStdout {
		command.Stdout = os.Stdout
	}
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

//创建审核
func createApplication(req apistructs.ApplicationCreateRequest) (*apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationCreateResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path("/api/applications").
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&req).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("create application error %s", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("create application not success %s", resp.Error.Msg)
	}

	if !r.IsOK() {
		return nil, fmt.Errorf("create application failed")
	}

	return &resp.Data, nil
}

func getApplication(appID uint64, userID string) (*apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationFetchResponse

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(conf.DiceOpenapiPublicUrl()).
		Path(fmt.Sprintf("/api/applications/%v", appID)).
		Header("Org-ID", strconv.FormatUint(conf.OrgId(), 10)).
		Header("USER-ID", userID).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return nil, fmt.Errorf("get application detail failed to request ("+err.Error()+")", false)
	}

	if !response.IsOK() {
		return nil, fmt.Errorf("get application detail failed to request, status-code: %d, content-type: %s", response.StatusCode(), response.ResponseHeader("Content-Type"))
	}

	if !resp.Success {
		return nil, fmt.Errorf("get application detailfailed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg)
	}

	return &resp.Data, nil
}

func getApplicationList() ([]apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationListResponse
	var b bytes.Buffer

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/applications")).
		Param("projectId", strconv.FormatUint(conf.ProjectId(), 10)).
		Param("pageSize", "9999").
		Param("pageNo", "1").
		Param("q", conf.ParamsApplicationName()).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s", response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()))
	}

	if !resp.Success {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg))
	}

	if resp.Data.Total == 0 {
		return nil, nil
	}

	return resp.Data.List, nil
}

func storeMetaFile(appID string, appExist bool) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "appId",
				Value: appID,
			},
			{
				Name:  "appExist",
				Value: fmt.Sprintf("%v", appExist),
			},
		},
	}

	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(conf.MetaFile(), string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func checkAppExist() (*apistructs.ApplicationDTO, error) {
	applications, err := getApplicationList()
	if err != nil {
		return nil, err
	}
	var existApp *apistructs.ApplicationDTO
	for _, app := range applications {
		if strings.EqualFold(app.Name, conf.ParamsApplicationName()) {
			existApp = &app
			break
		}
	}
	// application exist return appID and info
	if existApp != nil {
		return existApp, storeMetaFile(strconv.FormatUint(existApp.ID, 10), true)
	}
	return nil, nil
}
