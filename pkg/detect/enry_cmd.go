package detect

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/pkg/detect/bptype"
)

// detectDirLangsByEnryCmd return Language list ordered by percent desc already.
// Note: to use this function, you need install `enry` in your $PATH.
func detectDirLangsByEnryCmd(dir string) (Languages, error) {
	cmd := exec.Command("enry", "-prog", "-mode=byte", dir)
	cmd.Dir = dir
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// Original:
	// 98.16%	Go
	// 1.23%	Shell
	lines := strings.Split(string(b), "\n")
	var langs Languages
	for _, line := range lines {
		if len(line) <= 0 {
			continue
		}
		fs := strings.Fields(line)
		if len(fs) != 2 {
			return nil, errors.Errorf("cannot fields into two parts, strange line: %s", line)
		}
		f, err := strconv.ParseFloat(strings.TrimSuffix(fs[0], "%"), 64)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("get percent from: %s", fs[0]))
		}
		langs = append(langs, Language{Type: fs[1], Percent: f, Internal: bptype.IsInternalLang(fs[1])})
	}
	if !sort.IsSorted(langs) {
		sort.Sort(sort.Reverse(langs))
	}
	return langs, nil
}
