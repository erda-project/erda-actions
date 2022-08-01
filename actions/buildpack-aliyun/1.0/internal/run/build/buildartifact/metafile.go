package buildartifact

import (
	"encoding/json"
	"fmt"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/metadata"
	"github.com/erda-project/erda/pkg/strutil"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/pack"
)

func MetadataFromArtifact(artifact *apistructs.BuildArtifact) string {

	meta := make(metadata.Metadata, 0)

	// 镜像
	packResult := make([]pack.ModuleImage, 0)
	if err := json.Unmarshal([]byte(artifact.Content), &packResult); err == nil {
		for _, item := range packResult {
			meta = append(meta, metadata.MetadataField{
				Name: strutil.Concat(item.ModuleName, "-image"), Value: item.Image,
			})
		}
	}

	meta = append(meta,
		metadata.MetadataField{Name: "artifact.pipelineID", Value: fmt.Sprintf("%d", artifact.PipelineID)},
		metadata.MetadataField{Name: "artifact.type", Value: string(artifact.Type)},
		metadata.MetadataField{Name: "artifact.identityText", Value: artifact.IdentityText},
		metadata.MetadataField{Name: "artifact.content", Value: artifact.Content},
		metadata.MetadataField{Name: "artifact.clusterName", Value: artifact.ClusterName},
	)

	resp := Resp{Metadata: meta}
	respb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		bplog.Printf("failed to unmarshal metadata, err: %v\n", err)
		return "{}"
	}
	return string(respb)
}

type Resp struct {
	Metadata metadata.Metadata `json:"metadata"`
}
