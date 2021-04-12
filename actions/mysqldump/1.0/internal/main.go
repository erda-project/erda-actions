package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda/pkg/envconf"
)

type TableWithDumpOptions struct {
	Table                       string             `json:"table"`
	Where                       string             `json:"where"`
	DropTableIfExists           *bool              `json:"dropTableIfExists,omitempty"`           // if nil, use global config
	TryFilterByColumnsAndValues *map[string]string `json:"tryFilterByColumnsAndValues,omitempty"` // override global
	Charset                     string             `json:"charset,omitempty"`
}

type Conf struct {
	MySQLHost     string `env:"ACTION_MYSQL_HOST" required:"true"`
	MySQLPort     string `env:"ACTION_MYSQL_PORT" required:"true"`
	MySQLUsername string `env:"ACTION_MYSQL_USERNAME" required:"true"`
	MySQLPassword string `env:"ACTION_MYSQL_PASSWORD" required:"true"`
	MySQLDatabase string `env:"ACTION_MYSQL_DATABASE" required:"true"`

	GlobalDropTableIfExists   bool                    `env:"ACTION_GLOBAL_DROP_TABLE_IF_EXISTS" default:"true"`
	DumpAllTables             bool                    `env:"ACTION_DUMP_ALL_TABLES" default:"false"`
	DumpAllTablesIgnoreRegexp string                  `env:"ACTION_DUMP_ALL_TABLES_IGNORE_REGEXP"`
	MustIncludeTables         []*TableWithDumpOptions `env:"ACTION_MUST_INCLUDE_TABLES"`
	GlobalCharset             string                  `env:"ACTION_GLOBAL_CHARSET" default:"utf8mb4"`
	PostCommands              []string                `env:"ACTION_POST_COMMANDS"` // 可以使用 DUMP_FILE_PATH 获取 dump 文件

	// 若字段存在，则过滤；否则忽略
	GlobalTryFilterByColumnsAndValues map[string]string `env:"ACTION_GLOBAL_TRY_FILTER_BY_COLUMNS_AND_VALUES"`
	UploadDir                         string            `env:"UPLOADDIR" required:"true"`
}

type MySQLConn struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func getMySQLConnFromConfig(cfg Conf) *MySQLConn {
	return &MySQLConn{
		Host:     cfg.MySQLHost,
		Port:     cfg.MySQLPort,
		Username: cfg.MySQLUsername,
		Password: cfg.MySQLPassword,
		Database: cfg.MySQLDatabase,
	}
}

func (conn *MySQLConn) printMySQLConn() {
	logrus.Printf("\nMYSQL:\n"+
		"  host: %s\n"+
		"  port: %s\n"+
		"  user: %s\n"+
		"  db:   %s\n",
		conn.Host, conn.Port, conn.Username, conn.Database)
}

func main() {
	log.Init()

	var cfg Conf
	if err := envconf.Load(&cfg); err != nil {
		logrus.Fatalf("failed to load env config, err: %v", err)
	}
	var regex *regexp.Regexp
	if cfg.DumpAllTablesIgnoreRegexp != "" {
		_regex, err := regexp.Compile(cfg.DumpAllTablesIgnoreRegexp)
		if err != nil {
			logrus.Fatalf("invalid ignoreTablesRegexp: %s", regex)
		}
		regex = _regex
	}

	// conn
	conn := getMySQLConnFromConfig(cfg)

	conn.printMySQLConn()

	// one result file
	resultFilepath := filepath.Join("/tmp/dump_result/dump_all.sql")
	if _, err := os.OpenFile(resultFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		logrus.Fatalf("failed to open file for append dump result, err: %v", err)
	}

	os.Remove(resultFilepath)

	// handle tables
	// convert explicit tables to map
	explicitTablesMap := make(map[string]*TableWithDumpOptions)
	for _, tw := range cfg.MustIncludeTables {
		explicitTablesMap[tw.Table] = tw
	}
	var allTableAndWheres []*TableWithDumpOptions
	if cfg.DumpAllTables {
		allTables, err := conn.mysqlShowTables()
		if err != nil {
			logrus.Fatal(err)
		}
		includeTablesMap := make(map[string]struct{})
		for _, include := range cfg.MustIncludeTables {
			includeTablesMap[include.Table] = struct{}{}
		}
		for _, table := range allTables {
			// isInclude
			_, isInclude := includeTablesMap[table]
			// ignore
			if regex != nil {
				if regex.MatchString(table) && !isInclude {
					continue
				}
			}
			// already exist
			if existTw, ok := explicitTablesMap[table]; ok {
				allTableAndWheres = append(allTableAndWheres, existTw)
			} else {
				allTableAndWheres = append(allTableAndWheres, &TableWithDumpOptions{Table: table})
			}
		}
	} else {
		// only explicit tables
		allTableAndWheres = cfg.MustIncludeTables
	}
	for _, tw := range allTableAndWheres {
		// drop table config
		if tw.DropTableIfExists == nil {
			tw.DropTableIfExists = &cfg.GlobalDropTableIfExists
		}
		// filter by columns config only if where is empty
		if tw.Where == "" {
			if cfg.GlobalTryFilterByColumnsAndValues == nil {
				cfg.GlobalTryFilterByColumnsAndValues = make(map[string]string)
			}
			if tw.TryFilterByColumnsAndValues == nil {
				tw.TryFilterByColumnsAndValues = &cfg.GlobalTryFilterByColumnsAndValues
			}
			var queryColumns []string
			for column := range *tw.TryFilterByColumnsAndValues {
				queryColumns = append(queryColumns, column)
			}
			_, queryResult, err := conn.mysqlShowColumns(tw.Table, queryColumns...)
			if err != nil {
				logrus.Fatal(err)
			}
			for column, value := range cfg.GlobalTryFilterByColumnsAndValues {
				if queryResult[column] {
					tw.Where = addWhere(tw.Where, fmt.Sprintf("%s='%s'", column, value))
				}
			}
		}
		if tw.Charset == "" {
			tw.Charset = cfg.GlobalCharset
		}
		fmt.Println(tw.Table, " ", tw.Where)
	}

	for _, tw := range allTableAndWheres {
		if err := conn.mysqlDumpTable(*tw, resultFilepath); err != nil {
			logrus.Fatal(err)
		}
	}

	cmd := exec.Command("/bin/sh", "-c", "cp "+resultFilepath+" "+cfg.UploadDir)
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Fatal(err)
	}
	// post commands
	for _, postCmd := range cfg.PostCommands {
		logrus.Infof("begin run post command: %s", postCmd)
		cmd := exec.Command("/bin/sh", "-c", postCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		cmd.Env = append(cmd.Env, "DUMP_FILE_PATH="+resultFilepath)
		if err := cmd.Run(); err != nil {
			logrus.Fatal(err)
		}
	}
}

func addWhere(origin, new string) string {
	if origin == "" {
		return new
	}
	return fmt.Sprintf("%s AND %s", origin, new)
}
