package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (conn *MySQLConn) mysqlDumpTable(table TableWithDumpOptions, appendToFile string) error {
	args := []string{
		"mysqldump",
		"-h" + conn.Host,
		"-P" + conn.Port,
		"-u" + conn.Username,
		"-p" + conn.Password,
		// 不导出使用 GTID 方式开启了主从同步的 GTID 信息；
		// 若不关闭，则必须要求导入的目标数据库也需要开启 GTID
		"--set-gtid-purged=OFF",
		"--default-character-set=" + table.Charset,
	}
	if table.Where != "" {
		table.Where = strings.ReplaceAll(table.Where, `"`, `\"`)
		args = append(args, fmt.Sprintf(`--where="%s"`, table.Where))
	}
	if table.DropTableIfExists != nil && !*table.DropTableIfExists {
		args = append(args, "--no-create-info")
	}
	args = append(args, conn.Database, table.Table, ">>", appendToFile)
	fmt.Println(args)
	dump := exec.Command("/bin/sh", "-c", strings.Join(args, " "))
	dump.Stdout = os.Stdout
	dump.Stderr = os.Stdout
	if err := dump.Run(); err != nil {
		return fmt.Errorf("failed to dump table %s, err: %v", table.Table, err)
	}
	return nil
}
