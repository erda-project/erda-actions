package main

import (
	"fmt"
	"strconv"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

func getApplication(cfg *Conf) (*apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationFetchResponse

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(cfg.OpenAPIAddr).
		Path(fmt.Sprintf("/api/applications/%v", cfg.AppID)).
		Header("Org-ID", strconv.FormatUint(cfg.OrgId, 10)).
		Header("USER-ID", cfg.UserID).
		Header("Authorization", cfg.OpenAPIToken).Do().JSON(&resp)
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
