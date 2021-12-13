package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/release/1.0/internal/diceyml"
	"github.com/erda-project/erda-actions/pkg/docker"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
)

// migration flyway migration 文件镜像release
func migration(cfg *conf.Conf) (string, error) {
	if cfg.MigrationDir == "" {
		logrus.Info("empty migration dir.")
		return "", nil
	}
	if cfg.MigrationType == "erda" {
		return migrationErda(cfg)
	}
	if cfg.MigrationMysqlDatabase == "" {
		logrus.Info("migraion database must not be null.")
		return "", nil
	}
	logrus.Infof("migration dir: %v", cfg.MigrationDir)
	files, err := ioutil.ReadDir(cfg.MigrationDir)
	if err != nil {
		logrus.Errorf("read migration dir error: %v", err)
	}
	if len(files) == 0 {
		logrus.Info("there has nothing about migration sql.")
		return "", nil
	}
	for _, v := range files {
		logrus.Infof("migration file name is : %v", v.Name())
	}

	repo, err := packAndPushAppImage(cfg)
	if err != nil {
		logrus.Errorf("push migration image repo error: %v", err)
		return "", err
	}
	if repo == "" {
		return "", nil
	}
	// migration release dice.yml
	return migrationRelease(cfg, repo)
}

// migrationRelease migration信息release
func migrationRelease(cfg *conf.Conf, repo string) (string, error) {
	// generate release create request
	req := genReleaseRequest(cfg)
	req.ReleaseName = req.ReleaseName + "_migration"

	diceyml, err := genMigrationDiceYml(cfg, repo)
	if err != nil {
		logrus.Errorf("generate migration diceyml error: %v", err)
		return "", err
	}
	req.Dice = diceyml
	logrus.Infof("migration generate dice.yml: %v", req.Dice)
	// push release to dicehub
	releaseID, err := pushRelease(*cfg, req)
	if err != nil {
		return "", err
	}
	logrus.Infof("releaseId: %s", releaseID)
	// 切换工作目录
	logrus.Infof("switch back to the working directory: %s", cfg.WorkDir)
	if err := os.Chdir(cfg.WorkDir); err != nil {
		return "", err
	}
	return releaseID, nil
}

// genMigrationDiceYml 生成migration dice.yml
func genMigrationDiceYml(cfg *conf.Conf, repo string) (string, error) {
	diceymlContent, err := ioutil.ReadFile("dice.yml")
	if err != nil {
		return "", err
	}
	d, err := diceyml.New(diceymlContent)
	if err != nil {
		return "", errors.Wrap(err, "new parser failed")
	}
	images := map[string]string{"migration": repo}
	err = d.InsertImage(images)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert migration image to diceyml")
	}
	// 初始化mysql database环境变量
	diceymlCopy := d.Copy()
	diceymlCopy.SetEnv("MYSQL_DATABASE", cfg.MigrationMysqlDatabase)
	logrus.Infof("migration dice yml envs: %+v", diceymlCopy.Envs())

	yml, _ := yaml.Marshal(diceymlCopy.Obj())
	return string(yml), nil
}

// packAndPushAppImage 推送migration镜像，并返回镜像地址
func packAndPushAppImage(cfg *conf.Conf) (string, error) {
	// 切换工作目录
	if err := os.Chdir("/opt/action/comp/migration"); err != nil {
		return "", err
	}
	// copy assets
	if err := cp(cfg.MigrationDir, "./sql"); err != nil {
		return "", err
	}
	repo := getRepo(*cfg)

	if cfg.BuildkitEnable == "true" {
		if err := packWithBuildkit(repo); err != nil {
			return "", err
		}
	} else {
		if err := packWithDocker(repo); err != nil {
			return "", err
		}
	}

	// upload metadata
	if err := storeMigrationMetaFile(cfg, repo); err != nil {
		return "", err
	}
	fmt.Fprintf(os.Stdout, "successfully upload migration metafile\n")

	return repo, nil
}

func packWithDocker(repo string) error {
	packCmd := exec.Command("docker", "build",
		"-t", repo,
		"-f", "Dockerfile", ".")
	fmt.Fprintf(os.Stdout, "migration build, packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "migration build, successfully build app image: %s\n", repo)

	// docker push 业务镜像至集群 registry
	if err := docker.PushByCmd(repo, ""); err != nil {
		return err
	}
	return nil
}

func packWithBuildkit(repo string) error {
	packCmd := exec.Command("buildctl",
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=/opt/action/comp/migration",
		"--local", "dockerfile=/opt/action/comp/migration",
		"--output", "type=image,name=" + repo + ",push=true,registry.insecure=true")

	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "successfully build app image: %s\n", repo)
	return  nil
}

// storeMetaFile upload metadata
func storeMigrationMetaFile(cfg *conf.Conf, image string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "migration_image",
				Value: image,
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.Metafile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

// 生成业务镜像名称
func getRepo(cfg conf.Conf) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", "migration", time.Now().UnixNano())

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

func cp(a, b string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", a, b)
	cpCmd := exec.Command("cp", "-r", a, b)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
}
