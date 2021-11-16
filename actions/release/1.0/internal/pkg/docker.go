package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
)

func DockerBuildPushAndSetImages(cfg *conf.Conf) bool {

	if !ValidService(cfg.Services) {
		logrus.Errorf("check services value fail")
		return false
	}

	//根据多个service逐步构建
	for k, v := range cfg.Services {
		v.Name = k
		ok, image := buildDockerReturnImage(v, cfg)
		if !ok {
			return false
		}

		//因为是兼容images的，这里生成的imageName就塞入image就ok了
		if cfg.Images == nil {
			cfg.Images = map[string]string{}
		}
		cfg.Images[v.Name] = image
	}
	return true
}

//校验service的值
//todo 目前只是校验值不为空这些，是否可以实现校验值是否正确的校验
func ValidService(serviceMap map[string]conf.Service) bool {
	if serviceMap == nil {
		logrus.Errorf("services is empty")
		return false
	}

	for k, v := range serviceMap {
		if strings.Trim(k, " ") == "" {
			logrus.Infof("[release] services name %s is empty", k)
			return false
		}

		if strings.Trim(v.Cmd, " ") == "" {
			logrus.Infof("[release] services name %s cmd is empty", k)
			return false
		}

		if strings.Trim(v.Image, " ") == "" {
			logrus.Infof("[release] services name %s image is empty", k)
			return false
		}

		if v.Cps == nil || len(v.Cps) == 0 {
			logrus.Infof("[release] services name %s cps is empty, need cp your program to adaptation your run cmd", k)
			return false
		}

		for _, cp := range v.Cps {
			cpList := strings.Split(cp, ":")
			if len(cpList) != 2 {

				logrus.Infof("[release] services name %s cp is error, cp need you statement two Absolute path ,"+
					" like ${git-checkout}/target/run.jar:/root/run.jar and your cmd is java -jar /root/run.jar", k)
				return false
			}

			//for _, pathStr := range cpList {
			//	if path.IsAbs(pathStr) == false {
			//		logrus.Infof("[release] services name %s cp is error, cp need you statement two Absolute path ," +
			//			" like `${git-checkout}/target/run.jar:/root/run.jar` and your cmd is `java -jar /root/run.jar` ", k)
			//		return false
			//	}
			//}
		}
	}

	return true
}

func buildDockerReturnImage(service conf.Service, cfg *conf.Conf) (bool, string) {
	if !createDockerFile(service) {
		logrus.Errorf(" create dockerfile fail ")
		return false, ""
	}

	imageName := getRepoName(cfg, service.Name)
	logrus.Infof(" build docker image name %v ", imageName)
	if !buildDockerFile(service, imageName, cfg) {
		logrus.Errorf(" build dockerfile fail ")
		return false, ""
	}
	//
	//if !pushDockerFile(service, imageName) {
	//	logrus.Errorf(" push dockerfile fail ")
	//	return false, ""
	//}

	return true, imageName

}

//根据service和image的名称，用docker build构建其image
func buildDockerFile(service conf.Service, imageName string, cfg *conf.Conf) bool {
	dockerFileAddr := getDockerFileAddrByName(service.Name)

	//根据service种的cps，先将其需要拷贝到内容拷贝到当前目录
	//因为build是根据当前目录进行构建，要不然dockerfile的copy无法找到文件
	logrus.Infof("cps: %v", service.Cps)
	for _, cp := range service.Cps {
		cpList := strings.Split(cp, ":")
		f, _ := os.Stat(cpList[0])
		if f != nil && f.IsDir() {
			if err := runCommand(fmt.Sprintf("mkdir -p `pwd`/%s", getServiceTempPath(service.Name)+cpList[0])); err != nil {
				logrus.Errorf(" services %v: cp file before mkdir folder error: %v ", service.Name, err)
				return false
			}
			if err := runCommand(fmt.Sprintf("cp -r %s `pwd`/%s", cpList[0], getServiceTempPath(service.Name)+cpList[0])); err != nil {
				logrus.Errorf(" services %v: cp folder error: %v ", service.Name, err)
				return false
			}
		} else {
			fileNameDirArray := strings.Split(cpList[0], "/")
			fileName := fileNameDirArray[len(fileNameDirArray)-1]
			preDir := strings.Replace(cpList[0], fileName, "", 1)
			if err := runCommand(fmt.Sprintf("mkdir -p `pwd`/%s", getServiceTempPath(service.Name)+preDir)); err != nil {
				logrus.Errorf(" services %v: cp file before mkdir folder error: %v ", service.Name, err)
				return false
			}
			if err := runCommand(fmt.Sprintf("cp  %s `pwd`/%s", cpList[0], getServiceTempPath(service.Name)+cpList[0])); err != nil {
				logrus.Errorf(" services %v:  cp file error: %v ", service.Name, err)
				return false
			}
		}
	}

	//build
	if cfg.BuildkitEnable == "true" {
		return buildWithBuildkit(imageName, dockerFileAddr, service)
	} else {
		return buildWithDocker(imageName, dockerFileAddr, service)
	}
}


func buildWithDocker(imageName string, dockerFileAddr string, service conf.Service) bool {
	cmd := fmt.Sprintf(" docker build -t %s -f %s ./%s", imageName, dockerFileAddr, getServiceTempPath(service.Name))
	fmt.Println("build cmd :", cmd)
	if err := runCommand(cmd); err != nil {
		logrus.Errorf(" services %v:  docker run build error: %v ", service.Name, err)
		return false
	}

	if err := runCommand(fmt.Sprintf(" docker push %s", imageName)); err != nil {
		logrus.Errorf(" services %v:  docker run push error: %v ", service.Name, err)
		return false
	}
	return true
}

func buildWithBuildkit(imageName string, dockerFileAddr string, service conf.Service) bool {
	packCmd := exec.Command("buildctl",
		"--addr",
		"tcp://buildkitd.default.svc.cluster.local:1234",
		"--tlscacert=/.buildkit/ca.pem",
		"--tlscert=/.buildkit/cert.pem",
		"--tlskey=/.buildkit/key.pem",
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=./" + getServiceTempPath(service.Name),
		"--local", "dockerfile=" + filepath.Dir(dockerFileAddr),
		"--output", "type=image,name=" + imageName + ",push=true,registry.insecure=true")

	fmt.Fprintf(os.Stdout, "packCmd: %v\n", packCmd.Args)
	packCmd.Stdout = os.Stdout
	packCmd.Stderr = os.Stderr
	if err := packCmd.Run(); err != nil {
		logrus.Errorf(" services %v:  build image from buildkit error: %v ", service.Name, err)
		return false
	}
	return  true
}


func getServiceTempPath(serviceName string) string {
	return serviceName + "_temp/"
}

//根据服务名称和cfg中的一些值，构建出image的名称
func getRepoName(cfg *conf.Conf, serviceName string) string {
	repository := cfg.ProjectAppAbbr
	if repository == "" {
		repository = fmt.Sprintf("%s/%s", cfg.DiceOperatorId, random.String(8, random.Lowercase, random.Numeric))
	}
	tag := fmt.Sprintf("%s-%v", serviceName, time.Now().UnixNano())
	return strings.ToLower(fmt.Sprintf("%s/%s:%s", filepath.Clean(cfg.LocalRegistry), repository, tag))
}

//func pushDockerFile(service conf.Service, imageName string) bool {
//	if err := runCommand(fmt.Sprintf(" docker push %s", imageName)); err != nil {
//		logrus.Errorf(" services %v:  docker run push error: %v ", service.Name, err)
//		return false
//	}
//	return true
//}

//根据service中的值，构建对应的Dockerfile文件和内容
func createDockerFile(service conf.Service) bool {
	//写入image
	writeToDockerFile(fmt.Sprintf("FROM %s", service.Image), service.Name)
	for _, cp := range service.Cps {
		cpList := strings.Split(cp, ":")
		if len(cpList) != 2 {
			logrus.Infof("[release] services name %v cp is error, cp need you statement two Absolute path ,"+
				" like ${git-checkout}/target/run.jar:/root/run.jar and your cmd is java -jar /root/run.jar", service.Name)
			return false
		}
		writeToDockerFile(fmt.Sprintf("COPY %s %s", cpList[0], cpList[1]), service.Name)
	}
	writeToDockerFile(fmt.Sprintf("CMD %s ", service.Cmd), service.Name)
	return true
}

//根据服务名称获取其Dockerfile的绝对路径
func getDockerFileAddrByName(fileName string) string {
	return fmt.Sprintf("%s/Dockerfile", getDockerFileDirByName(fileName))
}

//根据服务名称获取其Dockerfile所属文件夹的绝对路径
func getDockerFileDirByName(fileName string) string {
	pwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s/dockerfiles_temp_2020_/%s", pwd, getServiceTempPath(fileName), fileName)
}

//根据服务名称，追加写入命令到对应的Dockerfile中
func writeToDockerFile(str string, serviceName string) {

	if !notExitCreateFile(serviceName) {
		return
	}

	runCommand(fmt.Sprintf("echo '%s' >> %s", str, getDockerFileAddrByName(serviceName)))
}

//根据服务名判断其Dockerfile文件是否存在，不存在就创建对应的文件
func notExitCreateFile(serviceName string) bool {
	if serviceName == "" {
		return false
	}

	filePath := getDockerFileAddrByName(serviceName)
	if filePath != "" && Exists(filePath) {
		return true
	}

	dir := getDockerFileDirByName(serviceName)
	if err := runCommand(fmt.Sprintf(" mkdir -p %s", dir)); err != nil {
		logrus.Infof("[release] create dir %s error: %v", dir, err)
		return false
	}

	if err := runCommand(fmt.Sprintf(" touch %s", filePath)); err != nil {
		logrus.Infof("[release] create file %s error: %v", filePath, err)
		return false
	}

	return true
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
