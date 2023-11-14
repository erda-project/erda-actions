package image

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/pkg/errors"
)

// ReportCacheImage 上报缓存镜像给 pipeline; 缓存镜像下载时，更新过期时间，延缓删除
func ReportCacheImage(openAPIAddr, openAPIToken string, reportReq apistructs.BuildCacheImageReportRequest) error {
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(openAPIAddr).
		Path("/api/build-caches").
		Header("Authorization", openAPIToken).
		JSONBody(&reportReq).
		Do().DiscardBody()
	if err != nil {
		return err
	}
	if !resp.IsOK() {
		return errors.Errorf("上报缓存镜像失败, status code: %d", resp.StatusCode())
	}

	return nil
}

func ReTagByGcrane(from, target string, insecure bool) error {
	fn := func(arg ...string) error {
		fmt.Fprintf(os.Stdout, "Run: gcrane, %v\n", arg)
		cmd := exec.Command("gcrane", arg...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	imageFile := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	args := []string{
		fmt.Sprintf("--insecure=%s", strconv.FormatBool(insecure)),
	}

	pullArgs := append([]string{"pull", from, imageFile}, args...)
	if err := fn(pullArgs...); err != nil {
		return err
	}

	pushArgs := append([]string{"push", imageFile, target}, args...)
	if err := fn(pushArgs...); err != nil {
		return err
	}

	return nil
}
