package main

import (
	"fmt"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"

	"github.com/erda-project/erda/pkg/envconf"
)

func TestAllTables(t *testing.T) {
	var cfg Conf
	envconf.MustLoad(&cfg)
	conn := getMySQLConnFromConfig(cfg)
	allTables, err := conn.mysqlShowTables()
	assert.NoError(t, err)
	fmt.Println(allTables)
}
