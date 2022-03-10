package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/project-package/1.0/internal/config"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

func Execute() error {
	logrus.SetOutput(os.Stdout)

	config, err := config.New()
	if err != nil {
		return err
	}

	fileId, err := exportPackage(config)
	if err != nil {
		logrus.Errorf("Export project package for %s failed, err: %v", config.ProjectName, err)
		return err
	}

	for i := 0; i < 12*config.WaitMinutes; i++ {
		record, err := getRecord(config, fileId)
		if err != nil {
			return err
		}
		logrus.Infof("Export project package for %s, state is %s", config.ProjectName, record.State)
		switch record.State {
		case apistructs.FileRecordStateFail:
			return errors.Errorf("Export project package for %s failed, err: %s", config.ProjectName, record.ErrorInfo)
		case apistructs.FileRecordStateSuccess:
			downloadUrl := fmt.Sprintf("%s/api/files/%s", config.DiceOpenapiPrefix, record.ApiFileUUID)
			storeMetaFile(config, downloadUrl)
			logrus.Infof("Download project package url: %s", downloadUrl)
			return nil
		case apistructs.FileRecordStatePending, apistructs.FileRecordStateProcessing:
			time.Sleep(5 * time.Second)
		}
	}

	logrus.Infof("Export project package for %s timeout, you may check state for record file %d in erda",
		config.ProjectName, fileId)

	return nil
}

func exportPackage(conf *config.Config) (uint64, error) {
	response := struct {
		apistructs.Header
		Data uint64
	}{}
	var b bytes.Buffer

	logrus.Infof("time: %s", time.Now().String())

	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/orgs/%d/projects/%d/package/actions/export", conf.OrgID, conf.ProjectID)).
		Header("Authorization", conf.CiOpenapiToken).
		Header("USER-ID", conf.UserID).
		Header("Org-ID", strconv.FormatUint(conf.OrgID, 10)).
		JSONBody(&conf.Artifacts).
		Do().Body(&b)

	if err != nil {
		return response.Data, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !resp.IsOK() {
		return response.Data, fmt.Errorf(
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				resp.StatusCode(), resp.ResponseHeader("Content-Type"), b.String()))
	}

	if err := json.Unmarshal(b.Bytes(), &response); err != nil {
		return response.Data, fmt.Errorf(
			fmt.Sprintf("failed to unmarshal project export response (" + err.Error() + ")"))
	}

	return response.Data, nil
}

func getRecord(conf *config.Config, id uint64) (apistructs.TestFileRecord, error) {
	var resp apistructs.GetTestFileRecordResponse
	var b bytes.Buffer

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/test-file-records/%d", id)).
		Header("Authorization", conf.CiOpenapiToken).
		Header("USER-ID", conf.UserID).
		Header("Org-ID", strconv.FormatUint(conf.OrgID, 10)).
		Do().Body(&b)
	if err != nil {
		return apistructs.TestFileRecord{}, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return apistructs.TestFileRecord{}, fmt.Errorf(
			fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s",
				response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()))
	}

	if err := json.Unmarshal(b.Bytes(), &resp); err != nil {
		return apistructs.TestFileRecord{}, fmt.Errorf(
			fmt.Sprintf("failed to unmarshal record response (" + err.Error() + ")"))
	}

	if !resp.Success {
		return apistructs.TestFileRecord{}, fmt.Errorf(
			fmt.Sprintf("failed to request, error code: %s, error message: %s",
				resp.Error.Code, resp.Error.Msg))
	}

	return resp.Data, nil
}

func storeMetaFile(conf *config.Config, downloadUrl string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "package_url",
				Value: downloadUrl,
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(conf.MetaFile, string(b), 0644); err != nil {
		return errors.New("write file:metafile failed")
	}
	return nil
}
