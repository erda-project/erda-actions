package pkg

import (
	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
)

// func TestValidService(t *testing.T) {

// 	tests := []struct {
// 		name   string
// 		image  string
// 		cmd    string
// 		cps    []string
// 		result bool
// 	}{
// 		{"", "", "", []string{""}, false},
// 		{" ", "", "", []string{""}, false},
// 		{"java-build", " ", "", []string{""}, false},
// 		{"java-build", "", "", []string{""}, false},
// 		{"java-build-1", "docker:v1.4", "", []string{""}, false},
// 		{"java-build-1", "docker:v1.4", " ", []string{""}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{""}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"aa/bb"}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"aa:bb:"}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"aa:bb"}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"/aa:bb"}, false},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"/aa:/bb"}, true},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"/aa/:/bb/as.jar"}, true},
// 		{"java-build-1", "docker:v1.4", "java -jar /root/as.jar", []string{"/aa/aa.jar:/bb/as.jar"}, true},
// 	}

// 	for _, v := range tests {

// 		services := make(map[string]conf.Service)
// 		services[v.name] = conf.Service{
// 			Cmd:   v.cmd,
// 			Image: v.image,
// 			Cps:   v.cps,
// 		}

// 		if ValidService(services) != v.result {
// 			t.Fail()
// 		}
// 	}

// }

// func TestDockerBuildPushAndSetImages(t *testing.T) {

// 	runCommand("rm -rf `pwd`/dockerfiles")

// 	if err := runCommand(fmt.Sprintf(" chmod -R 777 `pwd`")); err != nil {
// 		fmt.Sprintf("%v", err)
// 		t.Fail()
// 		return
// 	}

// 	runCommand(" echo `pwd`")

// 	tests := []struct {
// 		name   string
// 		image  string
// 		cmd    string
// 		cps    []string
// 		result bool
// 	}{
// 		//{"java-build-1", "docker", "ls", []string{"/Users/terminus/Music/网易云音乐/test.txt:/root/app"}, true},
// 		{"java-build-1", "docker", "ls", []string{"/Users/terminus/Music/网易云音乐/:/root/app"}, true},
// 		//{"java-build-2", "docker", "ls", []string{"/:~/"}, true},
// 		//{"java-build-3", "asdasda", "ls", []string{"/root/app:/root/app", "/root/app1:/root/app1"}, false},
// 	}

// 	for _, v := range tests {

// 		services := make(map[string]conf.Service)
// 		services[v.name] = conf.Service{
// 			Cmd:   v.cmd,
// 			Image: v.image,
// 			Cps:   v.cps,
// 			Name:  v.name,
// 		}

// 		cfg := buildCfg()
// 		cfg.Services = services
// 		if DockerBuildPushAndSetImages(cfg) != v.result {
// 			t.Fail()
// 		}

// 	}

// }

func buildCfg() *conf.Conf {
	cfg := conf.Conf{
		ProjectAppAbbr: "dice/release",
		DiceOperatorId: "",
		LocalRegistry:  "registry.erda.cloud",
	}
	return &cfg
}
