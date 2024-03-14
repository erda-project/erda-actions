package dockerfile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/erda-project/erda/pkg/strutil"
)

var (
	erdaUser = `RUN groupadd -g 1001 erda -f && useradd -u 1001 -g 1001 erda -o
	USER erda
	`
)

func InsertErdaUserToDockerfile(content []byte) []byte {
	lines := strutil.Split(string(content), "\n", true)
	var result []string
	var hasInserted bool
	for _, line := range lines {
		if strings.HasPrefix(line, "ENTRYPOINT") || strings.HasPrefix(line, "CMD") {
			result = append(result, erdaUser, line)
			hasInserted = true
			continue
		}
		result = append(result, line)
	}
	if !hasInserted {
		result = append(result, erdaUser)
	}
	return []byte(strings.Join(result, "\n"))
}

func ReplaceOrInsertBuildArgToDockerfile(content []byte, buildArgs map[string]string) []byte {

	// v 使用 json 序列化进行转义
	for k, v := range buildArgs {
		var vv bytes.Buffer
		d := json.NewEncoder(&vv)
		d.SetEscapeHTML(false)
		_ = d.Encode(v)
		buildArgs[k] = vv.String()
	}

	// polish Dockerfile
	lines := strutil.Split(string(content), "\n", true)

	// multiParts = global ARG + multi `FROM`
	var multiParts [][]string

	currentPart := make([]string, 0)
	for _, line := range lines {
		if !strings.HasPrefix(line, "FROM ") {
			currentPart = append(currentPart, line)
			continue
		}
		multiParts = append(multiParts, currentPart)
		currentPart = []string{line}
	}
	multiParts = append(multiParts, currentPart)

	// multiParts handle each part

	for i, part := range multiParts {
		lines := part

		copyBuildArgs := make(map[string]string, len(buildArgs))
		for k, v := range buildArgs {
			copyBuildArgs[k] = v
		}

		partResult := make([]string, 0, len(lines))

		// replace
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "ARG ") {
				line = line[4:]
				line = strings.TrimSpace(line)
				equalSignIndex := strings.Index(line, "=")
				spaceIndex := strings.Index(line, " ")
				// 用 = 分隔
				var splitSign = " "
				if spaceIndex == -1 { // 没有空格，那只能是 =
					splitSign = "="
				} else { //
					if equalSignIndex != -1 && equalSignIndex < spaceIndex {
						splitSign = "="
					}
				}
				sp := strings.SplitN(line, splitSign, 2)
				argK := sp[0]
				argV, ok := copyBuildArgs[argK]
				// 替换
				if ok {
					line = fmt.Sprintf(`ARG %s=%s`, argK, argV)
					delete(copyBuildArgs, argK)
				} else {
					line = fmt.Sprintf(`ARG %s`, line)
				}
			}
			partResult = append(partResult, line)
		}

		// buildArgs 字典序倒序，最终插入 Dockerfile 时为字典序正序
		orderedArgsSlice := make([]string, 0, len(copyBuildArgs))
		for k, v := range copyBuildArgs {
			orderedArgsSlice = append(orderedArgsSlice, fmt.Sprintf(`ARG %s=%s`, k, v))
		}
		sort.Sort(sort.StringSlice(orderedArgsSlice))

		// insert
		// 1. 只有 FROM 时无需插入 buildArg
		// 2. PART 第一行非 FROM，说明是 global arg，无需插入 buildArg
		if len(partResult) > 1 && strings.HasPrefix(partResult[0], "FROM ") {
			partResult = append([]string{partResult[0]}, append(orderedArgsSlice, partResult[1:]...)...)
		}

		multiParts[i] = partResult
	}

	// merge multiParts
	var result []string
	for _, part := range multiParts {
		for _, line := range part {
			result = append(result, line)
		}
		result = append(result, "")
	}
	return []byte(strings.Join(result, "\n"))
}

func TrimAllStringSpace(s string) string {
	reLeadcloseWhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := reLeadcloseWhtsp.ReplaceAllString(s, "")
	final = reInsideWhtsp.ReplaceAllString(final, " ")
	return final
}
