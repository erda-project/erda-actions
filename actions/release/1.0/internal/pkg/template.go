package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/release/1.0/internal/diceyml"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/cloudstorage"
)

// StorageURL 解析获取的存储URL
type StorageURL struct {
	Scheme   string
	UserName string
	Host     string
	Path     string
	PassWord string
}

// 存储类型
const (
	Oss  = "oss"
	File = "file"
)

func composeEnvYml(cfg *conf.Conf) (*diceyml.DiceYaml, error) {
	diceymlContent, err := ioutil.ReadFile(cfg.DiceYaml)
	if err != nil {
		return nil, err
	}
	d, err := diceyml.New(diceymlContent)
	if err != nil {
		return nil, errors.Wrap(err, "new parser failed")
	}

	switch cfg.Workspace {
	case string(apistructs.DevWorkspace):
		err = composeYaml(d, "development", cfg.DiceDevelopmentYaml)
	case string(apistructs.TestWorkspace):
		err = composeYaml(d, "test", cfg.DiceTestYaml)
	case string(apistructs.StagingWorkspace):
		err = composeYaml(d, "staging", cfg.DiceStagingYaml)
	case string(apistructs.ProdWorkspace):
		err = composeYaml(d, "production", cfg.DiceProductionYaml)
	default:
		return nil, errors.Errorf("invalid workspace")
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to compose diceyml")
	}

	return d, nil
}

func composeYaml(targetYml *diceyml.DiceYaml, env, envYmlFile string) error {
	if _, err := os.Stat(envYmlFile); os.IsNotExist(err) {
		return nil
	}

	envContent, err := ioutil.ReadFile(envYmlFile)
	if err != nil {
		return err
	}
	envYml, err := diceyml.New(envContent)
	if err != nil {
		return err
	}

	err = targetYml.Compose(env, envYml)
	if err != nil {
		return err
	}

	return nil
}

func useRetagImage(d *diceyml.DiceYaml, cfg *conf.Conf) {
	if d == nil {
		return
	}

	diceServices := d.Services()
	if diceServices == nil {
		return
	}
	for k, service := range cfg.Services {
		if diceServices[k] == nil {
			continue
		}
		// 若RetagImage没有被定义过，则跳过此逻辑
		if service.RetagImage == "" {
			continue
		}
		// UseRetagImage 默认为nil(或值为true)，使用retag中的镜像；在pipeline.yml中可以使用 useRetagImage: false 将值修改成false
		if service.UseRetagImage == nil || *service.UseRetagImage {
			diceServices[k]["image"] = service.RetagImage
		}
	}
}

func insertCommands(d *diceyml.DiceYaml, cfg *conf.Conf) {
	if d == nil {
		return
	}

	services := cfg.Services
	diceServices := d.Services()
	if diceServices == nil {
		return
	}
	for k, v := range services {
		if diceServices[k] == nil {
			continue
		}

		diceServices[k]["cmd"] = v.Cmd
	}
	d.SetServices(diceServices)
	return
}

func insertImages(d *diceyml.DiceYaml, cfg *conf.Conf) error {
	images := make(map[string]string, len(cfg.ReplacementImages))
	for _, v := range cfg.ReplacementImages {
		err := parseImages(v, images)
		if err != nil {
			return err
		}
	}
	if len(cfg.Images) > 0 { // 允许用户指定service对应外部镜像，不用打包的镜像，且以用户指定镜像优先
		for k, v := range cfg.Images {
			images[k] = v
		}
	}

	err := d.InsertImage(images)
	if err != nil {
		return errors.Wrap(err, "failed to insert image to diceyml")
	}

	return nil
}

func executeSQL(cfg *conf.Conf, storage *StorageURL, d *diceyml.DiceYaml) error {
	if _, err := os.Stat(cfg.InitSQL); os.IsNotExist(err) {
		return err
	}

	sqlURL, err := dBSQLPush(cfg.WorkDir, cfg.InitSQL, storage)
	if err != nil {
		return errors.Errorf("failed to store init sql, (%+v)", err)
	}

	if err = insertAddons(d, cfg.Workspace, "mysql", map[string]string{"init_sql": sqlURL}); err != nil {
		return errors.Errorf("failed to insert initsql into dice.yml, (%+v)", err)
	}

	return nil
}

func insertAddons(d *diceyml.DiceYaml, env, addonType string, context map[string]string) error {
	if err := d.InsertAddonOptions(transfer(env), addonType, context); err != nil {
		return err
	}
	return nil
}

func parseImages(f string, images map[string]string) error {
	fileValue, err := ioutil.ReadFile(f) //TODO file path correct?
	if err != nil {
		return errors.Wrapf(err, "read file %s failed", f)
	}

	kv := []struct {
		ModuleName string `json:"module_name"`
		Image      string `json:"image"`
	}{}
	err = json.Unmarshal(fileValue, &kv)
	if err != nil {
		return errors.Wrapf(err, "json %s unmarshal failed", string(fileValue))
	}

	for _, item := range kv {
		images[item.ModuleName] = item.Image
	}
	return nil
}

func transfer(env string) diceyml.EnvType {
	switch env {
	case "DEV":
		return diceyml.DevEnv
	case "TEST":
		return diceyml.TestEnv
	case "STAGING":
		return diceyml.StagingEnv
	case "PROD":
		return diceyml.ProdEnv
	default:
		return diceyml.BaseEnv
	}
}

func dBSQLPush(wd, sourceDir string, storage *StorageURL) (string, error) {
	var URL string

	if err := tarDBSql(wd, sourceDir); err != nil {
		return "", err
	}

	dbPath := filepath.Join(wd, "db.tar.gz")
	ossPath := fmt.Sprintf("terminus-initdb/%s/db.tar.gz", strconv.FormatInt(time.Now().UnixNano(), 10))

	URL, err := StorageFile(dbPath, ossPath, storage)
	if err != nil {
		return "", err
	}

	logrus.Info("successed to push sql to storage")

	return URL, nil
}

// UploadFiles 上传用户的文件列表
func UploadFiles(files string, storage *StorageURL) ([]apistructs.ReleaseResource, error) {
	fileList := strings.Split(files, ",")

	releaseResources := []apistructs.ReleaseResource{}
	for _, file := range fileList {
		ossPath := fmt.Sprintf("terminus-files/%s/%s", strconv.FormatInt(time.Now().UnixNano(), 10), file)
		URL, err := StorageFile(strings.Trim(file, " "), ossPath, storage)
		if err != nil {
			return nil, err
		}

		if URL == "" {
			logrus.Warningf("upload file with return nil url, file: %s", file)
			continue
		}

		// insert release resource
		resource := apistructs.ReleaseResource{
			Type: checkResourceType(file),
			Name: file,
			URL:  URL,
		}
		releaseResources = append(releaseResources, resource)
	}

	return releaseResources, nil
}

func checkResourceType(file string) apistructs.ResourceType {
	if strings.Contains(file, string(apistructs.ResourceTypeDiceYml)) {
		return apistructs.ResourceTypeDiceYml
	}

	if strings.Contains(file, string(apistructs.ResourceTypeAddonYml)) {
		return apistructs.ResourceTypeAddonYml
	}

	if strings.Contains(file, string(apistructs.ResourceTypeSQL)) {
		return apistructs.ResourceTypeSQL
	}

	if strings.Contains(file, string(apistructs.ResourceTypeDiceYml)) {
		return apistructs.ResourceTypeDiceYml
	}

	return apistructs.ResourceTypeDataSet
}

// StorageFile 将文件上传到oss或者网盘存储
func StorageFile(filePath, ossPath string, storage *StorageURL) (string, error) {
	var URL string
	switch storage.Scheme {
	case File:
		timeStamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		storagePath := filepath.Join(storage.Path, timeStamp)
		command := fmt.Sprintf("mkdir -p %s; cp %s %s", storagePath, filePath, storagePath)
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}

		if strings.Contains(filePath, "db.tar.gz") {
			URL = fmt.Sprint("file://", filepath.Join(timeStamp, "db.tar.gz"))
		} else {
			URL = fmt.Sprint("file://", storagePath)
		}
	case Oss:
		ossURL, err := Push2CloudStorage(storage, filePath, ossPath)
		if err != nil || ossURL == "" {
			return "", errors.Errorf("failed to upload file to cloud storage, (%+v)", err)
		}
		URL = ossURL
	default:
		return "", errors.Errorf("error type, storage: %s", storage)
	}

	return URL, nil
}

func tarDBSql(wd, sourceDir string) error {
	script := strings.Join([]string{
		fmt.Sprintf("#!/bin/sh\n"),
		fmt.Sprintf("tar -czPf %s %s\n", filepath.Join(wd, "db.tar.gz"), sourceDir),
	}, "")

	scriptPath := filepath.Join(wd, "upload_db.sh")
	if err := CreateFile(scriptPath, script, 0755); err != nil {
		return errors.Errorf("failed to create upload_db.sh, (%+v)", err)
	}

	if err := ExecScript(scriptPath); err != nil {
		return errors.Errorf("failed to exec upload_db.sh, (%+v)", err)
	}

	return nil
}

// Push2CloudStorage 存储到对象存储
func Push2CloudStorage(storageURL *StorageURL, dbPath, ossPath string) (string, error) {
	var (
		err    error
		url    string
		client cloudstorage.Client
	)

	if client, err = cloudstorage.New(fmt.Sprint("http://", storageURL.Host),
		storageURL.UserName, storageURL.PassWord); err != nil {
		return "", err
	}

	if url, err = client.UploadFile(getBucket(strings.TrimPrefix(storageURL.Path, "/")),
		ossPath, dbPath); err != nil {
		return "", err
	}

	return url, nil
}

// CreateFile 创建文件
func CreateFile(absPath, content string, perm os.FileMode) error {
	if !filepath.IsAbs(absPath) {
		return errors.Errorf("not an absolute path: %s", absPath)
	}
	err := os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return errors.Wrap(err, "make parent dir error")
	}
	f, err := os.OpenFile(filepath.Clean(absPath), os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "write content to file error")
	}
	return nil
}

// ExecScript 执行command
func ExecScript(scriptPath string) error {
	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// 默认为 terminus-dice
func getBucket(bucket string) string {
	if bucket != "" {
		return bucket
	}
	return "terminus-dice"
}
