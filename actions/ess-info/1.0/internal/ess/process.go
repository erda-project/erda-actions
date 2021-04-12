package ess

import (
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	hostsFile       = "hosts"
	instanceIDsFile = "instance_ids"
)

func (e *Ess) writeFile(filePath string, content []string) error {
	err := filehelper.CreateFile(filePath, strings.Join(content, "\n"), 0755)
	if err != nil {
		logrus.Errorf("create file failed, file path: %s, error: %s", filePath, err)
	}
	return nil
}

func Process() error {
	var privateIPs []string
	var instanceIDs []string
	e := Ess{}
	envconf.MustLoad(&e)

	result, err := e.GetEssInfo()
	if err != nil {
		return err
	}
	if result != nil {
		for id, ip := range result {
			privateIPs = append(privateIPs, ip)
			instanceIDs = append(instanceIDs, id)
		}
	}

	// write metadata to WORKDIR, and it will be used by next pipeline stages
	hosts := strings.Join(privateIPs, ",")
	hostsFilePath := filepath.Join(e.WorkDir, hostsFile)
	if err := e.writeFile(hostsFilePath, []string{hosts}); err != nil {
		return err
	}

	ids := strings.Join(instanceIDs, ",")
	instanceIDsFilePath := filepath.Join(e.WorkDir, instanceIDsFile)
	if err := e.writeFile(instanceIDsFilePath, []string{ids}); err != nil {
		return err
	}
	return nil
}
