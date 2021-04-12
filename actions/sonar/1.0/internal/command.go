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
	resultMap, err := command.packageManager.Analysis(cfg)
	if len(resultMap) > 0 {
		logrus.Println("Sonar analysis result:")
		for k, v := range resultMap {
			fmt.Printf("%s: %v\n", k, v)
		}
		if err := storeMetaFile(cfg.MetaFile, resultMap); err != nil {
			logrus.Fatalf("failed to store meta file, err: %v", err)
		}
	}
	return err
}

func storeMetaFile(metafilepath string, result map[ResultKey]string) error {
	var kvs []string
	for k, v := range result {
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
	}
	content := strutil.Join(kvs, "\n", true)
	return filehelper.CreateFile(metafilepath, content, 0644)
}
