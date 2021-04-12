package buildartifact

import (
	"encoding/json"
	"fmt"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/pack"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/strutil"
)

func MetadataFromArtifact(artifact *apistructs.BuildArtifact) string {

	metadata := make(apistructs.Metadata, 0)

	// 镜像
	packResult := make([]pack.ModuleImage, 0)
	if err := json.Unmarshal([]byte(artifact.Content), &packResult); err == nil {
		for _, item := range packResult {
			metadata = append(metadata, apistructs.MetadataField{
				Name: strutil.Concat(item.ModuleName, "-image"), Value: item.Image,
			})
		}
	}

	metadata = append(metadata,
		apistructs.MetadataField{Name: "artifact.pipelineID", Value: fmt.Sprintf("%d", artifact.PipelineID)},
		apistructs.MetadataField{Name: "artifact.type", Value: string(artifact.Type)},
		apistructs.MetadataField{Name: "artifact.identityText", Value: artifact.IdentityText},
		apistructs.MetadataField{Name: "artifact.content", Value: artifact.Content},
		apistructs.MetadataField{Name: "artifact.clusterName", Value: artifact.ClusterName},
	)

	resp := Resp{Metadata: metadata}
	respb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		bplog.Printf("failed to unmarshal metadata, err: %v\n", err)
		return "{}"
	}
	return string(respb)
}

type Resp struct {
	Metadata apistructs.Metadata `json:"metadata"`
}
