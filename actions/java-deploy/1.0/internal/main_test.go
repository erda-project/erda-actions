package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/conf"
	"github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/run/dlog"
)

func TestJavaDeploy_Monolith(t *testing.T) {
	os.Args[0] = os.Getenv("GOPATH") + "/src/github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/main.go"
	os.Setenv("ACTION_WORKDIR", filepath.Join(filepath.Dir(os.Args[0]), "testdata/monolith"))
	os.Setenv("ACTION_REGISTRY", "http://localhost:8081/repository/maven-snapshots/")
	os.Setenv("ACTION_USERNAME", "admin")
	os.Setenv("ACTION_PASSWORD", "admin123")

	os.Setenv("METAFILE", "/tmp/123/metafile")
	os.Setenv("WORKDIR", ".")
	os.Setenv("BP_NEXUS_URL", "https://repo.terminus.io")
	os.Setenv("BP_NEXUS_USERNAME", "readonly")
	os.Setenv("BP_NEXUS_PASSWORD", "Hello1234")

	/////////////////////////////////////
	// main
	/////////////////////////////////////
	logrus.SetOutput(os.Stdout)

	err := conf.LoadEnvConfig()
	if err != nil {
		dlog.Fatalf("failed to load env config, err: %v\n", err)
	}

	err = run.Execute()
	if err != nil {
		dlog.Fatal(err)
	}
}

func TestJavaDeploy_Cmd(t *testing.T) {
	os.Args[0] = os.Getenv("GOPATH") + "/src/github.com/erda-project/erda-actions/actions/java-deploy/1.0/internal/main.go"
	os.Setenv("ACTION_WORKDIR", filepath.Join(filepath.Dir(os.Args[0]), "testdata/monolith"))
	os.Setenv("ACTION_REGISTRY", "http://localhost:8081/repository/maven-snapshots/")
	os.Setenv("ACTION_USERNAME", "admin")
	os.Setenv("ACTION_PASSWORD", "admin123")
	os.Setenv("ACTION_CMD", "./gradlew publish")

	os.Setenv("METAFILE", "/tmp/123/metafile")
	os.Setenv("WORKDIR", ".")
	os.Setenv("BP_NEXUS_URL", "https://repo.terminus.io")
	os.Setenv("BP_NEXUS_USERNAME", "readonly")
	os.Setenv("BP_NEXUS_PASSWORD", "Hello1234")

	/////////////////////////////////////
	// main
	/////////////////////////////////////
	logrus.SetOutput(os.Stdout)

	err := conf.LoadEnvConfig()
	if err != nil {
		dlog.Fatalf("failed to load env config, err: %v\n", err)
	}

	err = run.Execute()
	if err != nil {
		dlog.Fatal(err)
	}
}
