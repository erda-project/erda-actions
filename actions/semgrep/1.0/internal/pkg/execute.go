package pkg

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/sirupsen/logrus"
)

var outputFile = "/tmp/result.sarif"

func (s *Semgrep) Execute() error {
	if s.cfg.CodeDir == "" {
		return fmt.Errorf("missing context")
	}
	if s.cfg.Config == "" {
		return fmt.Errorf("missing config rule")
	}
	s.cmd.SetDir(s.cfg.CodeDir)
	s.cmd.Add("--config")
	s.cmd.Add(s.cfg.Config)
	if s.cfg.Format != "" {
		s.cmd.Add(fmt.Sprintf("--%s", s.cfg.Format))
	} else {
		s.cmd.Add("--sarif")
	}
	s.cmd.Add("-o")
	s.cmd.Add(outputFile)
	for _, arg := range s.cfg.Args {
		logrus.Infof("semgrep add arg: %s", arg)
		s.cmd.Add(arg)
	}
	defer func() {
		if err := s.results.Store(); err != nil {
			logrus.Errorf("failed to store results: %v", err)
		}
		return
	}()
	logrus.Infof("running semgrep: %s", s.cmd.String())
	if err := s.cmd.Run(); err != nil {
		return fmt.Errorf("semgrep ci failed, err: %v", err)
	}
	uploadFile, err := UploadFileNew(outputFile, s.cfg)
	if err != nil {
		return fmt.Errorf("failed to upload file, err: %v", err)
	}
	s.results.Add("download_url", uploadFile.Data.DownloadURL)
	return nil
}

// UploadFileNew 上传到api/files接口
func UploadFileNew(filePath string, cfg *Conf) (*apistructs.FileDownloadFailResponse, error) {
	logrus.Infof("upload file %s", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileName := filepath.Base(filePath)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile %v", err)
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}

	var resp apistructs.FileDownloadFailResponse
	request := httpclient.New(httpclient.WithCompleteRedirect()).Post(cfg.PlatformParams.OpenAPIAddr).
		Path("/api/files").
		Param("fileFrom", "release").
		Param("public", "true").
		Header("Content-Type", writer.FormDataContentType()).
		Header("Authorization", cfg.PlatformParams.OpenAPIToken)
	httpResp, err := request.RawBody(body).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, fmt.Errorf("failed to upload file, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Error.Msg)
	}
	return &resp, nil
}
