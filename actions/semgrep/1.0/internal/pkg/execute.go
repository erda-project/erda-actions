package pkg

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda-proto-go/core/file/pb"
	"github.com/erda-project/erda/pkg/filehelper"
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
	args := []string{"semgrep"}
	for _, arg := range s.cfg.Args {
		logrus.Infof("semgrep add arg: %s", arg)
		args = append(args, arg)
	}
	args = append(args, fmt.Sprintf("--config=%s", s.cfg.Config))
	if s.cfg.Format != "" {
		args = append(args, fmt.Sprintf("--%s", s.cfg.Format))
	} else {
		args = append(args, "--sarif")
	}
	args = append(args, "-o", outputFile)
	defer func() {
		if err := s.results.Store(); err != nil {
			logrus.Errorf("failed to store results: %v", err)
		}
		return
	}()
	logrus.Infof("running semgrep args: %s", strings.Join(args, " "))
	os.Chdir(s.cfg.CodeDir)
	cmd := exec.Command("/bin/sh", "-c", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = s.cfg.CodeDir
	cmd.Env = NewEnv()
	cmd.Run()
	if err := filehelper.CheckExist(outputFile, false); err != nil {
		return fmt.Errorf("semgrep ci execute failed, err: %v", err)
	}
	uploadFile, err := UploadFileNew(outputFile, s.cfg)
	if err != nil {
		return fmt.Errorf("failed to upload file, err: %v", err)
	}
	s.results.Add("download_url", uploadFile.Data.DownloadURL)
	return nil
}

func NewEnv() []string {
	env := []string{
		"PATH=/opt/go/bin:/go/bin:/opt/nodejs/bin:/opt/maven/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	}
	return env
}

// UploadFileNew 上传到api/files接口
func UploadFileNew(filePath string, cfg *Conf) (*pb.FileUploadResponse, error) {
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

	var resp pb.FileUploadResponse
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
