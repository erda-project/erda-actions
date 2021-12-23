package build

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/ios/1.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/dice"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/retry"
)

const (
	localTarFilePath        = "/tmp/context.tar"
	resultTmpPath           = "/tmp/result"
	MaxQueryTaskFailedCount = 5
)

var cfg conf.Conf

type RunnerTask struct {
	TaskID         string `json:"task_id"`
	Status         string `json:"status"` // pending running success failed
	ContextDataUrl string `json:"context_data_url"`
	ResultDataUrl  string `json:"result_data_url"`
}

func Execute() error {

	envconf.MustLoad(&cfg)
	os.Chdir(cfg.Context)

	commands := []string{}
	if cfg.P12Cert != nil {
		if path.Ext(cfg.P12Cert.Dest) != ".p12" {
			return errors.New("invalid p12 ext")
		}

		logrus.Infof("copy p12 file")
		destDir := filepath.Dir(cfg.P12Cert.Dest)
		err := os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			logrus.Warnf("error create p12 cert dir %s", destDir)
		}
		err = cp(cfg.P12Cert.Source, cfg.P12Cert.Dest)
		if err != nil {
			return err
		}
		commands = append(commands,
			fmt.Sprintf("security import %s -k  {{osx_keychains_path}}  -P %v -A", cfg.P12Cert.Dest, cfg.P12Cert.Password))
	}
	if cfg.MobileProvision != nil {
		if path.Ext(cfg.MobileProvision.Dest) != ".mobileprovision" {
			return errors.New("invalid provision ext")
		}
		logrus.Infof("copy provision file")
		destDir := filepath.Dir(cfg.MobileProvision.Dest)
		err := os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			logrus.Warnf("error create provision dir %s", destDir)
		}
		err = cp(cfg.MobileProvision.Source, cfg.MobileProvision.Dest)
		if err != nil {
			return err
		}
		commands = append(commands, fmt.Sprintf("open %s", cfg.MobileProvision.Dest))
	}
	// 生成移动应用打包配置文件
	if err := createIOSBuildCfg(); err != nil {
		return err
	}
	// 打包本地context环境
	logrus.Infof("package local context tar")
	os.Chdir(cfg.PipelineContext)
	err := runCommand("rm ", "-rf", ".git")
	if err != nil {
		fmt.Printf("remove .git dir error %v", err)
	}
	err = Tar(localTarFilePath, ".")
	if err != nil {
		return err
	}
	uploadReq := &dice.UploadFileRequest{
		FilePath:      localTarFilePath,
		OpenApiPrefix: cfg.DiceOpenapiPrefix,
		Token:         cfg.CiOpenapiToken,
		From:          "runner-action",
		Public:        "true",
		ExpireIn:      "3600s",
	}
	var uploadResult *apistructs.FileUploadResponse
	err = retry.DoWithInterval(
		func() error {
			var err error
			uploadResult, err = dice.UploadFile(uploadReq)
			return err
		}, 10, time.Second*15,
	)
	if err != nil {
		return err
	}

	commands = append(commands, cfg.Commands...)
	createReq := &apistructs.CreateRunnerTaskRequest{
		JobID:          cfg.PipelineTaskLogID,
		ContextDataUrl: uploadResult.Data.DownloadURL,
		Commands:       commands,
		WorkDir:        strings.Replace(cfg.Context, cfg.PipelineContext+"/", "", -1),
		Targets:        cfg.Targets,
	}

	taskID, err := CreateTask(cfg, createReq)
	if err != nil {
		return err
	}

	go ListenSignal(cfg, taskID)
	queryTaskResult, err := WaitTask(cfg, strconv.FormatInt(taskID, 10))
	if err != nil {
		return err
	}

	logrus.Infof("download result")

	localTargetTarPath := "/tmp/target.tar"
	err = retry.DoWithInterval(
		func() error {
			os.RemoveAll(localTarFilePath)
			return dice.DownloadFile(queryTaskResult.ResultDataUrl, localTargetTarPath)
		}, 10, time.Second*15,
	)
	if err != nil {
		return err
	}

	logrus.Infof("extract remote result to workDir")

	err = os.MkdirAll(resultTmpPath, os.ModePerm)
	if err != nil {
		return err
	}
	// 解压到本地 workDir
	err = UnTar(localTargetTarPath, resultTmpPath)
	if err != nil {
		return err
	}

	for _, target := range cfg.Targets {
		cp(path.Join(resultTmpPath, target), cfg.WorkDir)
	}

	fmt.Fprintln(os.Stdout, "target files")
	runCommand("ls ", cfg.WorkDir)

	return nil

}

func cp(srcDir, destDir string) error {
	fmt.Fprintf(os.Stdout, "copying %s to %s\n", srcDir, destDir)
	cpCmd := exec.Command("cp", "-r", srcDir, destDir)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	return cpCmd.Run()
}

func runCommand(cmd ...string) error {
	c := strings.Join(cmd, " ")
	command := exec.Command("/bin/bash", "-c", c)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return err
	}
	return nil
}

func WaitTask(cfg conf.Conf, taskID string) (*apistructs.RunnerTask, error) {
	queryFailedCount := 0
	isRunning := false
	for true {
		time.Sleep(time.Second * 10)
		task, err := QueryTask(cfg, taskID)
		if err != nil {
			logrus.Errorf("failed to query task status err:%s", err)
			queryFailedCount += 1
			if queryFailedCount > MaxQueryTaskFailedCount {
				return nil, err
			}
			continue
		}
		queryFailedCount = 0
		switch task.Status {
		case "pending":
			logrus.Infof("pending to process")
		case "failed":
			return nil, fmt.Errorf("failed to build")
		case "success":
			return task, nil
		case "running":
			if !isRunning {
				isRunning = true
				logrus.Infof("start running")
			}
		default:
			logrus.Infof("task status: %s", task.Status)
		}
	}
	return nil, errors.New("unknown task status")
}

func ListenSignal(cfg conf.Conf, taskID int64) {
	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		select {
		case sig := <-sigChan:
			logrus.Infof("received a signal: %s (%d)", sig, sig)

			switch sig {
			case syscall.SIGTERM:
				logrus.Infof("cancel job")
				err := UpdateTaskStatus(cfg, taskID, "canceled")
				if err != nil {
					logrus.Errorf("failed to cancel job err:%s", err)
				}
				os.Exit(int(syscall.SIGTERM))

			default:
				break
			}
		}
	}
}

func createIOSBuildCfg() error {
	return filehelper.CreateFile(cfg.Context+"/mobileBuild.cfg", "PIPELINE_ID="+cfg.PipelineID, 0644)
}
