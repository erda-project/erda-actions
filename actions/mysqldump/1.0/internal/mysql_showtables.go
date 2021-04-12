package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (conn *MySQLConn) mysqlShowTables() ([]string, error) {
	cmd := exec.Command("mysql",
		"-h"+conn.Host,
		"-P"+conn.Port,
		"-u"+conn.Username,
		"-p"+conn.Password,
		conn.Database,
		"--default-character-set=utf8mb4",
		"-e", `show tables\G`,
	)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to show tables, err: %v", err)
	}
	var allTables []string
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	const prefix = "Tables_in_"
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		ss := strings.SplitN(line, " ", 2)
		if len(ss) != 2 {
			continue
		}
		allTables = append(allTables, ss[1])
	}
	return allTables, nil
}
