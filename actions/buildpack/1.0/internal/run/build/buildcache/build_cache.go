package buildcache

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/util"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

func ReportCacheImage(action string) {
	var err error
	defer func() {
		if err != nil {
			bplog.Printf("上报缓存镜像失败，请忽略。镜像：%s，失败原因：%v\n", conf.EasyUse().CalculatedCacheImage, err)
		}
	}()

	addr, err := util.GetDiceOpenAPIAddress()
	if err != nil {
		return
	}

	registerReq := apistructs.BuildCacheImageReportRequest{
		Action:      action,
		ClusterName: conf.PlatformEnvs().ClusterName,
		Name:        conf.EasyUse().CalculatedCacheImage,
	}

	var result apistructs.BuildCacheImageReportResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(addr).Path("/api/build-caches").
		Header("Authorization", conf.PlatformEnvs().OpenAPIToken).
		JSONBody(&registerReq).Do().JSON(&result)
	if err != nil {
		return
	}
	if !r.IsOK() {
		err = errors.Errorf("status-code %d, result %v", r.StatusCode(), result)
		return
	}

	bplog.Println("上报缓存镜像成功!")
}
