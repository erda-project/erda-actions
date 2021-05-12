package dlog

import (
	"strings"

	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/bplog"
	"github.com/erda-project/erda-actions/actions/buildpack-aliyun/1.0/internal/run/langdetect/types"
)

const (
	contextRoot = "context root"
)

// find file
func TryToFindFileUnderDir(tryFile, underDir string) {
	bplog.Printf("try to find file %q under %s\n", tryFile, underDir)
}
func TryToFindFileInPath(fileName, filePath string) {
	bplog.Printf("try to find file %q in %s", fileName, handleFilePath(filePath))
}
func TryToFindFileUnderContextRoot(tryFile string) {
	TryToFindFileUnderDir(tryFile, contextRoot)
}
func FindFileInPath(fileName, filePath string) {
	bplog.Printf("find %q in %s\n", fileName, handleFilePath(filePath))
}
func FindFileUnderContextRoot(fileName string) {
	FindFileInPath(fileName, contextRoot)
}
func NotFoundFileInPath(fileName, filePath string) {
	bplog.Printf("not found %q in %s\n", fileName, handleFilePath(filePath))
}
func NotFoundFileUnderContextRoot(fileName string) {
	NotFoundFileInPath(fileName, contextRoot)
}

// find content
func TryToFindContentInFile(content, filePath string) {
	bplog.Printf("try to find content %q inside file %s\n", content, handleFilePath(filePath))
}
func FindContentInFile(content, filePath string) {
	bplog.Printf("find %q in file: %s\n", content, handleFilePath(filePath))
}
func NotFoundContentInFile(content, filePath string) {
	bplog.Printf("not found content %q inside file %s\n", content, handleFilePath(filePath))
}

// expect
func TryToDetectExpectedBuildType(expect types.BuildType) {
	bplog.Printf("try to detect expected build_type: %s\n", expect)
}
func TryToDetectExpectedContainerType(expect types.ContainerType) {
	bplog.Printf("try to detect expected container_type: %s\n", expect)
}

func handleFilePath(path string) string {
	actionContextDirPrefix := "/.pipeline/container/context/"
	if strings.HasPrefix(path, actionContextDirPrefix) {
		last := strings.TrimPrefix(path, actionContextDirPrefix)
		slashIndex := strings.Index(last, "/")
		if slashIndex != -1 {
			last = "${" + last[:slashIndex] + "}" + last[slashIndex:]
		} else {
			last = "${" + last + "}"
		}
		return last
	}
	return path
}
