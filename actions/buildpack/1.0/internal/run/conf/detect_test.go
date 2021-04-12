package conf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack/1.0/internal/run/langdetect/types"
)

func TestDetectLang(t *testing.T) {
	type TestCase struct {
		caseName      string
		cfg           conf
		language      types.Language
		buildType     types.BuildType
		containerType types.ContainerType
	}
	testcases := []TestCase{
		// java
		{
			caseName: "compatible: only set bp_repo java",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termjava.git",
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMaven,
			containerType: types.ContainerTypeSpringBoot,
		},
		{
			caseName: "compatible: bp_repo: java, bp_ver: release/3.10",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termjava.git",
					BpVer:  "release/3.10",
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMaven,
			containerType: types.ContainerTypeSpringBoot,
		},
		{
			caseName: "compatible: bp_repo: java, bp_ver: feature/edas-dubbo-1.0.8",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termjava.git",
					BpVer:  "feature/edas-dubbo-1.0.8",
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMavenEdas,
			containerType: types.ContainerTypeEdas,
		},
		{
			caseName: "dir: java monolith, lang: java",
			cfg: conf{
				params: &params{
					Language: types.LanguageJava,
					Context:  "../../../bp/java/build/maven/testdata/monolith",
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMaven,
			containerType: types.ContainerTypeSpringBoot,
		},
		{
			caseName: "dir: java monolith, lang: java, build_type: maven",
			cfg: conf{
				params: &params{
					Language:  types.LanguageJava,
					Context:   "../../../bp/java/build/maven/testdata/monolith",
					BuildType: types.BuildTypeMaven,
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMaven,
			containerType: types.ContainerTypeSpringBoot,
		},
		{
			caseName: "dir: java monolith, lang: java, build_type: maven_edas",
			cfg: conf{
				params: &params{
					Language:  types.LanguageJava,
					Context:   "../../../bp/java/build/maven/testdata/monolith",
					BuildType: types.BuildTypeMavenEdas,
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMavenEdas,
			containerType: types.ContainerTypeEdas,
		},
		{
			caseName: "dir: java monolith, lang: java, build_type: maven_edas, container_type: edas",
			cfg: conf{
				params: &params{
					Language:      types.LanguageJava,
					Context:       "../../../bp/java/build/maven/testdata/monolith",
					BuildType:     types.BuildTypeMavenEdas,
					ContainerType: types.ContainerTypeEdas,
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMavenEdas,
			containerType: types.ContainerTypeEdas,
		},
		{
			caseName: "dir: java multi-modules, lang: java",
			cfg: conf{
				params: &params{
					Language: types.LanguageJava,
					Context:  "../../../bp/java/build/maven/testdata/multi-modules",
				},
			},
			language:      types.LanguageJava,
			buildType:     types.BuildTypeMaven,
			containerType: types.ContainerTypeSpringBoot,
		},
		// node
		{
			caseName: "compatible: only set bp_repo nodejs",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termnodejs.git",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeHerd,
		},
		{
			caseName: "compatible: only set bp_repo spa",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termspa.git",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeSpa,
		},
		{
			caseName: "compatible: bp_repo: node, bp_ver: release/3.10",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termnodejs.git",
					BpVer:  "release/3.10",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeHerd,
		},
		{
			caseName: "compatible: bp_repo: node, bp_ver: master",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termnodejs.git",
					BpVer:  "master",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeHerd,
		},
		{
			caseName: "compatible: bp_repo: spa, bp_ver: release/3.10",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termspa.git",
					BpVer:  "release/3.10",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeSpa,
		},
		{
			caseName: "compatible: bp_repo: spa, bp_ver: master",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termspa.git",
					BpVer:  "master",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeSpa,
		},
		{
			caseName: "dir: testdata/herd",
			cfg: conf{
				params: &params{
					Context: "../../../bp/node/build/npm/testdata/herd",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeHerd,
		},
		{
			caseName: "dir: testdata/spa",
			cfg: conf{
				params: &params{
					Context: "../../../bp/node/build/npm/testdata/spa",
				},
			},
			language:      types.LanguageNode,
			buildType:     types.BuildTypeNpm,
			containerType: types.ContainerTypeSpa,
		},
		// dockerfile
		{
			caseName: "compatible: only set bp_repo dockerfile",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termdockerimage.git",
				},
			},
			language:      types.LanguageDockerfile,
			buildType:     types.BuildTypeDockerfile,
			containerType: types.ContainerTypeDockerfile,
		},
		{
			caseName: "compatible: bp_repo: dockerfile, bp_ver: release/3.10",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termdockerimage.git",
					BpVer:  "release/3.10",
				},
			},
			language:      types.LanguageDockerfile,
			buildType:     types.BuildTypeDockerfile,
			containerType: types.ContainerTypeDockerfile,
		},
		{
			caseName: "compatible: bp_repo: dockerfile, bp_ver: master",
			cfg: conf{
				params: &params{
					BpRepo: "file:///opt/action/buildpacks/dice-bpack-termdockerimage.git",
					BpVer:  "master",
				},
			},
			language:      types.LanguageDockerfile,
			buildType:     types.BuildTypeDockerfile,
			containerType: types.ContainerTypeDockerfile,
		},
	}

	for _, tc := range testcases {
		bplog.Printf("test case: %s\n", tc.caseName)
		err := detectLang(tc.cfg.params)
		require.NoError(t, err)
		require.Equal(t, tc.language, tc.cfg.params.Language)
		require.Equal(t, tc.buildType, tc.cfg.params.BuildType)
		require.Equal(t, tc.containerType, tc.cfg.params.ContainerType)
		fmt.Println()
	}
}
