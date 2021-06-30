package buildartifact

import (
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/util"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/http/httpclientutil"
)

func QueryBuildArtifact(sha string) (*apistructs.BuildArtifact, error) {
	addr, err := util.GetDiceOpenAPIAddress()
	if err != nil {
		return nil, err
	}

	var artifact apistructs.BuildArtifact
	r := httpclient.New(httpclient.WithCompleteRedirect()).Get(addr).Path("/api/build-artifacts/"+sha).
		Header("Authorization", conf.PlatformEnvs().OpenAPIToken)
	if err = httpclientutil.DoJson(r, &artifact); err != nil {
		return nil, err
	}
	return &artifact, nil
}

func RegisterBuildArtifact(artifactSHA string, identityText string, packResult []byte) (*apistructs.BuildArtifact, error) {
	addr, err := util.GetDiceOpenAPIAddress()
	if err != nil {
		return nil, err
	}

	registerReq := apistructs.BuildArtifactRegisterRequest{
		SHA:          artifactSHA,
		Type:         string(apistructs.BuildArtifactOfFileContent),
		Content:      string(packResult),
		PipelineID:   conf.PlatformEnvs().PipelineID,
		ClusterName:  conf.PlatformEnvs().ClusterName,
		IdentityText: identityText,
	}

	var artifact apistructs.BuildArtifact
	r := httpclient.New(httpclient.WithCompleteRedirect()).Post(addr).Path("/api/build-artifacts").
		Header("Authorization", conf.PlatformEnvs().OpenAPIToken).
		JSONBody(&registerReq)
	if err = httpclientutil.DoJson(r, &artifact); err != nil {
		return nil, err
	}
	return &artifact, nil
}
