package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// mysqlShowColumns return allColumnNames, columnQueryResultMap, error
func (conn *MySQLConn) mysqlShowColumns(table string, queryColumns ...string) ([]string, map[string]bool, error) {
	cmd := exec.Command("mysql",
		"-h"+conn.Host,
		"-P"+conn.Port,
		"-u"+conn.Username,
		"-p"+conn.Password,
		conn.Database,
		"--default-character-set=utf8mb4",
		"-e", fmt.Sprintf("show columns from `%s` \\G", table),
	)
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to show columns, table: %s, err: %v", table, err)
	}
	queryColumnMap := make(map[string]struct{})
	for _, column := range queryColumns {
		queryColumnMap[column] = struct{}{}
	}
	var allColumnNames []string
	columnQueryResultMap := make(map[string]bool)
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	const prefix = "Field:"
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		ss := strings.SplitN(line, " ", 2)
		if len(ss) != 2 {
			continue
		}
		columnName := ss[1]
		allColumnNames = append(allColumnNames, columnName)
		// 该字段是否需要查询
		if _, ok := queryColumnMap[columnName]; ok {
			columnQueryResultMap[columnName] = true
		}
	}
	// 若请求字段不存在，则赋值 false
	for _, query := range queryColumns {
		if _, ok := columnQueryResultMap[query]; !ok {
			columnQueryResultMap[query] = false
		}
	}
	return allColumnNames, columnQueryResultMap, nil
}
