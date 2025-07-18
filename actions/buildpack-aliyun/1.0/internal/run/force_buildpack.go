package run

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/build/buildartifact"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
)

// handleForceBuildpack need force if artifact is nil
func handleForceBuildpack(artifactSHA, identityText string) (artifact *apistructs.BuildArtifact) {

	// 如果用户指定了强制打包，则执行正常流程
	if conf.Params().ForceBuildpack {
		bplog.Println("用户指定了强制打包!")
		return
	}

	// 根据打包上下文，由系统判定是否需要强制打包，如果判定为需要(例如 java: grep snapshots pom.xml)，则自动强制打包
	force := systemJudgeForceBuildpack()

	bplog.Printf("本次构建的 SHA 计算基于:\n%s\n", identityText)
	bplog.Printf("本次构建的 SHA: %s\n", artifactSHA)

	// 如果无需强制打包，则尝试查询 Artifact
	if !force {
		bplog.Printf("开始查询历史构建产物 (SHA: %s) ......\n", artifactSHA)
		artifact, err := buildartifact.QueryBuildArtifact(artifactSHA)
		if err != nil {
			bplog.Println("查询历史构建产物失败，将执行正常打包流程")
			return nil
		}
		bplog.Println("查询历史构建产物成功!")
		switch artifact.Type {
		case string(apistructs.BuildArtifactOfFileContent):
			if err := filehelper.CreateFile(filepath.Join(conf.PlatformEnvs().WorkDir, "pack-result"), artifact.Content, 0644); err != nil {
				bplog.Printf("根据历史构建产物创建 pack-result 失败，将执行正常打包流程。失败原因: %v\n", err)
				return nil
			}
			bplog.Println("使用历史构建产物，无需重新打包!")
			return artifact
		case string(apistructs.BuildArtifactOfNfsLink):
			if err := filehelper.Copy(artifact.Content, filepath.Join(conf.PlatformEnvs().WorkDir, "pack-result")); err != nil {
				bplog.Printf("根据历史构建产物拷贝 pack-result 失败，将执行正常打包流程。失败原因: %v\n", err)
				return nil
			}
			bplog.Println("使用历史构建产物，无需重新打包!")
			return artifact
		default:
			bplog.Printf("构建产物类型 (%s) 不支持，请忽略\n", artifact.Type)
		}
	}

	return nil
}

func systemJudgeForceBuildpack() bool {
	bplog.Println("系统开始判定是否需要强制打包 ......")
	var forceBuildpack = false
	if conf.Params().Language == types.LanguageJava {
		forceBuildpack = systemJudgeForceBuildpackJava()
	}
	if forceBuildpack {
		bplog.Println("系统判定需要强制打包!")
	} else {
		bplog.Println("系统判定无需强制打包!")
	}
	return forceBuildpack
}

// systemJudgeForceBuildpackJava grep all pom.xml has snapshots or not
func systemJudgeForceBuildpackJava() (forceBuildpack bool) {
	var err error
	defer func() {
		if err != nil {
			bplog.Printf("系统判定是否需要强制打包失败，自动设置为强制打包。失败原因: %v\n", err)
			forceBuildpack = true
		}
	}()

	var script = []string{"#!/bin/bash", "set -eo pipefail"}
	pomDir := filepath.Join(conf.PlatformEnvs().WorkDir, ".cache_pom")
	script = append(script,
		"#!/bin/bash",

		"cd "+conf.Params().Context,

		"pom_dir="+pomDir,
		"mkdir -p ${pom_dir}",

		// parent pom
		`echo ">>> copy parent pom file ......"`,
		"parent_pom_dir=${pom_dir}/parent_pom",
		"mkdir -p ${parent_pom_dir} && cp -f pom.xml ${parent_pom_dir}/pom.xml",

		// all pom
		`echo ">>> copy all pom files ......"`,
		`all_pom_dir=${pom_dir}/all_pom`,
		"mkdir -p ${all_pom_dir}",
		`find . -name 'pom.xml' -exec cp --parents {} ${all_pom_dir} \;`,

		// list pom
		`echo ">>> list pom.xml in ${pom_dir}:"`,
		`find ${pom_dir} -type f -name 'pom.xml' | cut -d '/' -f 6-`,

		// grep snapshots
		`echo ">>> grep snapshots in all pom.xml"`,
		`set +o pipefail`,
		`any_snapshot=$(grep -Rn -- '-SNAPSHOT<' ${pom_dir} | grep -v -- '<!--' | grep -v -- '-->' | grep -- 'version>' | cut -d '/' -f 7- | tee snapshots.txt | wc -l)`,
		`set -o pipefail`,
		`if [ ${any_snapshot} -gt 0 ]; then`,
		`cat snapshots.txt`,
		`echo ">>> maven dependency 中含有 SNAPSHOT jar, 为了保证正确性, 系统自动触发强制打包!"`,
		`exit 1`, // 退出码非0，则为需要强制打包
		`else`,
		`echo ">>> maven dependency 中没有找到 SNAPSHOT jar, 无需强制打包!"`,
		`exit 0`, // 退出码为0，则不需要强制打包
		`fi`,
	)

	scriptPath := filepath.Join(conf.PlatformEnvs().WorkDir, "whether_query_artifact.sh")
	if err := filehelper.CreateFile(scriptPath, strings.Join(script, "\n"), 0755); err != nil {
		return true
	}
	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Dir = filepath.Dir(conf.PlatformEnvs().WorkDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return true
	}
	return false
}
