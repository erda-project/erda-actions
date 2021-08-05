package cancel

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/app-run/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

func Cancel() error {
	if err := conf.Load(); err != nil {
		return err
	}

	pipelineIDStr, err := getPipelineInfo()
	if err != nil {
		return err
	}

	pipelineID, err := strconv.ParseUint(pipelineIDStr, 10, 64)
	if err != nil {
		return err
	}

	return cancelPipeline(pipelineID)
}

func cancelPipeline(pipelineID uint64) error {
	var resp apistructs.PipelineCreateResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/cicds/%v/actions/cancel", conf.PipelineID())).
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&apistructs.PipelineCancelRequest{PipelineID: pipelineID}).
		Do().JSON(&resp)

	if err != nil {
		return fmt.Errorf("cancel pipeline error %s", err)
	}

	if !resp.Success {
		return fmt.Errorf("cancel pipeline not success %s", resp.Error.Msg)
	}

	if !r.IsOK() {
		return fmt.Errorf("cancel pipeline failed")
	}

	return nil
}

func getPipelineInfo() (string, error) {
	fileValue, err := ioutil.ReadFile(filepath.Join(conf.WorkDir(), "pipelineInfo"))
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file pipelineInfo")
	}

	if fileValue == nil {
		return "", errors.New("null pipelineInfo content")
	}

	pipelineIDInfo := strings.Split(string(fileValue), "=")
	if len(pipelineIDInfo) != 2 {
		return "", errors.New("failed to get pipelineID")
	}

	return pipelineIDInfo[1], nil
}
