package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

// UploadFileNew 上传到api/files接口
func UploadFileNew(filePath string, cfg conf.Conf) (*apistructs.FileDownloadFailResponse, error) {
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
	request := httpclient.New(httpclient.WithCompleteRedirect()).Post(cfg.DiceOpenapiPrefix).
		Path("/api/files").
		Param("fileFrom", "release").
		Param("public", "true").
		Header("Content-Type", writer.FormDataContentType()).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.RawBody(body).Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, errors.Errorf("failed to upload file, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, errors.Errorf(resp.Error.Msg)
	}

	// 兖矿临时修改
	v := strings.ReplaceAll(resp.Data.DownloadURL, "http://dice.paas.ykjt.cc", "https://appstore.ykjt.cn")
	v = strings.ReplaceAll(v, "https://dice.paas.ykjt.cc", "https://appstore.ykjt.cn")
	resp.Data.DownloadURL = v
	return &resp, nil
}

func Zip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
