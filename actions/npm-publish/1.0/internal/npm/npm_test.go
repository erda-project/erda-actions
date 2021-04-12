package npm

import (
	"testing"
)

func Test_Login(t *testing.T) {
	npm := NewNPM()

	err := npm.Login("deyang.ddy",
		"ddy805786",
		"deyang.ddy@alibaba-inc.com",
		"https://registry.npmjs.org",
	)
	if err != nil {
		t.Error("Npm login err:%v", err)
	} else {
		t.Log("Npm login OK!")
	}
}

func Test_Publish(t *testing.T) {
	npm := NewNPM()

	err := npm.Publish("/Users/ddy/ddy-publish-test", "ddy-test", "https://registry.npmjs.org")
	if err != nil {
		t.Error("Npm publish err:%v", err)
	} else {
		t.Log("Npm publish OK!")
	}
}

func Test_View(t *testing.T) {
	npm := NewNPM()
	var packageInfo *PackageInfo

	packageInfo, err := npm.View("ddy-publish-test", "")
	if err != nil {
		t.Error("Npm view err:%s", err)
	} else {
		t.Log("Npm view OK!\n packageInfo:\n%v\n", *packageInfo)
	}
}

func Test_Logout(t *testing.T) {
	npm := NewNPM()

	err := npm.Logout("https://registry.npmjs.org")
	if err != nil {
		t.Error("Npm logout err:%v", err)
	} else {
		t.Log("Npm logout OK!")
	}
}

func Test_Install(t *testing.T) {
	npm := NewNPM()

	err := npm.Install("npm-cli-login", "https://registry.npmjs.org", true)
	if err != nil {
		t.Error("Npm install err:%v", err)
	} else {
		t.Log("Npm install OK!")
	}
}
