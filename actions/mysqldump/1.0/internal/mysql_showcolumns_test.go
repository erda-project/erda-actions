package main

import (
	"fmt"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"github.com/erda-project/erda/pkg/envconf"
)

func TestShowColumns(t *testing.T) {
	var cfg Conf
	envconf.MustLoad(&cfg)
	conn := getMySQLConnFromConfig(cfg)
	allColumns, columnQueryResultMap, err := conn.mysqlShowColumns("ps_orgs", "org_id", "cluster_name", "id")
	assert.NoError(t, err)
	fmt.Println(allColumns)
	fmt.Println(columnQueryResultMap)
}
