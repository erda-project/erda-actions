package run

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/dlog"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/pom"
	"github.com/erda-project/erda-actions/pkg/render"
)

func executeMvn() (string, error) {
	// create settings.xml
	settingsXMLPath := filepath.Join(filepath.Dir(os.Args[0]), "assets", "settings.xml")
	if err := render.RenderTemplate(filepath.Dir(settingsXMLPath), map[string]string{
		"DEPLOY_USERNAME": conf.UserConf().Username,
		"DEPLOY_PASSWORD": conf.UserConf().Password,

		"CLUSTER_NEXUS_USERNAME": conf.PlatformConf().ClusterNexusUsername,
		"CLUSTER_NEXUS_PASSWORD": conf.PlatformConf().ClusterNexusPassword,
		"CLUSTER_NEXUS_URL":      conf.PlatformConf().ClusterNexusURL,
	}); err != nil {
		return "", errors.Errorf("failed to render settings.xml, err: %v\n", err)
	}

	// execute `mvn clean deploy`
	mvnDeployArgs := []string{
		"clean",
		"deploy",
		"-e", "-B", "-U",
		"-DaltDeploymentRepository=deploy::default::" + conf.UserConf().Registry,
		"--settings=" + settingsXMLPath,
	}
	if conf.UserConf().SkipTests {
		mvnDeployArgs = append(mvnDeployArgs, "-Dmaven.test.skip=true")
	}
	if conf.UserConf().Modules != "" {
		mvnDeployArgs = append(mvnDeployArgs, "-am", "-pl", conf.UserConf().Modules)
	}

	deployCmd := exec.Command("mvn", mvnDeployArgs...)
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	deployCmd.Dir = conf.UserConf().Workdir
	dlog.Printf("will execute cmd: %s\n", deployCmd.String())

	if err := deployCmd.Run(); err != nil {
		return "", errors.Errorf("failed to execute java deploy cmd, err: %v", err)
	}

	// 写入 metafile
	var metaContent string
	// gav
	gav, err := pom.GetGAV(filepath.Join(conf.UserConf().Workdir, "pom.xml"))
	if err == nil {
		metaContent += fmt.Sprintf("%s=%s\n", pom.GroupID, gav.GroupID)
		metaContent += fmt.Sprintf("%s=%s\n", pom.ArtifactID, gav.ArtifactID)
		metaContent += fmt.Sprintf("%s=%s\n", pom.Version, gav.Version)
		dlog.Printf("GAV:\n%s\n", metaContent)
	}

	return metaContent, nil
}
