package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)

type Command struct {
	packageManager PackageManager
}

func NewCommand(packageManager PackageManager) *Command {
	return &Command{
		packageManager: packageManager,
	}
}

func (command *Command) Analysis(cfg *Conf) error {
	results, err := command.packageManager.Analysis(cfg)
	if results != nil && len(*results) > 0 {
		if err := storeMetaFile(cfg.MetaFile, results); err != nil {
			logrus.Fatalf("failed to store meta file, err: %v", err)
		}
	}
	return err
}

func storeMetaFile(metafilepath string, results *ResultMetas) error {
	var kvs []string
	for _, meta := range *results {
		kvs = append(kvs, fmt.Sprintf("%s=%s", meta.Key, meta.Value))
	}
	content := strutil.Join(kvs, "\n", true)
	return filehelper.CreateFile(metafilepath, content, 0644)
}
