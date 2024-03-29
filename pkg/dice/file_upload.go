package dice

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-proto-go/core/file/pb"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type UploadFileRequest struct {
	FilePath      string
	OpenApiPrefix string
	Token         string
	From          string
	Public        string
	ExpireIn      string
}

// UploadFile 上传到api/files接口
func UploadFile(req *UploadFileRequest, timeout time.Duration) (*pb.FileUploadResponse, error) {
	logrus.Infof("upload file %s", req.FilePath)
	f, err := os.Open(req.FilePath)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(req.FilePath)
	multiparts := map[string]httpclient.MultipartItem{
		"file": {
			Reader:   f,
			Filename: fileName,
		},
	}
	var resp pb.FileUploadResponse
	request := httpclient.New(httpclient.WithCompleteRedirect(), httpclient.WithTimeout(time.Second, timeout)).Post(req.OpenApiPrefix).
		Path("/api/files").
		Param("fileFrom", req.From).
		Param("expiredIn", req.ExpireIn).
		Param("public", req.Public).
		Header("Authorization", req.Token).
		MultipartFormDataBody(multiparts)
	httpResp, err := request.Do().JSON(&resp)
	if err != nil {
		logrus.Errorf("err request %s", err)
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, errors.Errorf("failed to upload file, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, errors.Errorf(resp.Error.Msg)
	}
	return &resp, nil
}
