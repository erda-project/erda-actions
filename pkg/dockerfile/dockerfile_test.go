package dockerfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/erda-project/erda/pkg/filehelper"

	"gopkg.in/stretchr/testify.v1/require"
)

func TestTrimAllStringSpace(t *testing.T) {
	s := ` ARG DEP_CMD="hello  world"  `
	fmt.Println(TrimAllStringSpace(s))
}

func TestReplaceOrInsertBuildArgToDockerfile(t *testing.T) {
	dockerfile := []byte(`
   

FROM alpine:3.7

ARG DEP_CMD="npm ci"

RUN echo ${DEP_CMD}

RUN eval ${DEP_CMD}
`)

	bpArgs := map[string]string{
		"DEP_CMD":               `echo hello && echo I'm linjun && echo I"m linjun`,
		"FORCE_DEP":             "true",
		"FORCE_UPDATE_SNAPSHOT": "false",
	}

	b := ReplaceOrInsertBuildArgToDockerfile(dockerfile, bpArgs)
	fmt.Println(string(b))
	fmt.Println()
	fmt.Println()

	dockerfilePath, err := filepath.Abs(filepath.Join("./testdata/Dockerfile"))
	require.NoError(t, err)
	dockerfileContent, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)
	newDockerfileContent := ReplaceOrInsertBuildArgToDockerfile(dockerfileContent, bpArgs)
	fmt.Println(string(newDockerfileContent))
	err = filehelper.CreateFile(dockerfilePath, string(newDockerfileContent), 0644)
	require.NoError(t, err)
}

func TestReplaceOrInsertBuildArgToDockerfile2(t *testing.T) {

	bpArgs := map[string]string{
		"DEP_CMD":               `echo hello && echo I'm linjun && echo I"m linjun`,
		"FORCE_DEP":             "true",
		"FORCE_UPDATE_SNAPSHOT": "false",
	}

	bpArgs = nil

	dockerfilePath, err := filepath.Abs(filepath.Join("./testdata/Dockerfile"))
	require.NoError(t, err)
	dockerfileContent, err := ioutil.ReadFile(dockerfilePath)
	fmt.Println(len(dockerfileContent))
	require.NoError(t, err)
	newDockerfileContent := ReplaceOrInsertBuildArgToDockerfile(dockerfileContent, bpArgs)
	fmt.Println(string(newDockerfileContent))
	fmt.Println(len(newDockerfileContent))
	err = filehelper.CreateFile(dockerfilePath, string(newDockerfileContent), 0644)
	require.NoError(t, err)
}

func TestReplaceOrInsertBuildArgToDockerfile3(t *testing.T) {
	bpArgs := map[string]string{
		"DATE":    "20200114",
		"DEP_CMD": `cat /etc/hosts && sed "s/127.0.0.1/localhost/g" /etc/hosts && echo done && echo '<>'`,
		"DD":      "DD",
		"A":       "B",
	}
	dockerfilePath, err := filepath.Abs(filepath.Join("./testdata/Dockerfile.multi-stage"))
	require.NoError(t, err)
	dockerfileContent, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)
	newDockerfileContent := ReplaceOrInsertBuildArgToDockerfile(dockerfileContent, bpArgs)
	fmt.Println(string(newDockerfileContent))
}

func TestReplaceOrInsertBuildArgToDockerfile4(t *testing.T) {
	dockerfile := []byte(`
FROM alpine:3.7

ARG URL
`)
	bpArgs := map[string]string{
		"BACKSLASH":   `\`,
		"SINGLEQUOTE": `'`,
		"DOUBLEQUOTE": `"`,
		"SLASH":       `/`,
	}
	fmt.Println(string(ReplaceOrInsertBuildArgToDockerfile(dockerfile, bpArgs)))
}

func TestInsertErdaUserToDockerfile(t *testing.T) {
	dockerfile := []byte(`
FROM registry.erda.cloud/retag/pyroscope-java:v0.11.5 as pyroscope-java
FROM registry.erda.cloud/erda-x/openjdk:8_11

ARG CONTAINER_VERSION=v8
ENV CONTAINER_VERSION ${CONTAINER_VERSION}

ENV SCRIPT_ARGS ${SCRIPT_ARGS}

COPY comp/openjdk/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY pre_start.sh /pre_start.sh
RUN chmod +x /pre_start.sh

COPY comp/fonts /usr/share/fonts/custom
#COPY comp/arthas-boot.jar /
COPY comp/jacocoagent.jar /opt/jacoco/jacocoagent.jar

ARG ERDA_VERSION
COPY comp/spot-agent/${ERDA_VERSION}/spot-agent.tar.gz /tmp/spot-agent.tar.gz
RUN \
	if [ "${MONITOR_AGENT}" = true ]; then \
        mkdir -p /opt/spot; tar -xzf /tmp/spot-agent.tar.gz -C /opt/spot; \
	fi && rm -rf /tmp/spot-agent.tar.gz

ENTRYPOINT ["/entrypoint.sh"]
`)
	fmt.Println(string(InsertErdaUserToDockerfile(dockerfile)))
}
