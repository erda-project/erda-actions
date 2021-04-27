package util

import (
	"fmt"
	"sort"
)

func GetSortedKeySlice(m map[string]string) []string {
	if m == nil {
		return nil
	}
	var keys []string
	for k, v := range m {
		keys = append(keys, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(keys)
	return keys
}
