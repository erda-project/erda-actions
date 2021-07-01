package image

import (
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
