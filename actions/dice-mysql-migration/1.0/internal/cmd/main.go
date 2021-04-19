package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/erda-project/erda-actions/actions/dice-mysql-migration/1.0/internal/log"
	"github.com/erda-project/erda-actions/actions/dice-mysql-migration/1.0/internal/migration"
	"github.com/erda-project/erda/pkg/sqlparser/migrator"
)

func main() {
	go startSandbox()

	log.Infoln("Erda MySQL Migration start working")
	log.Infof("Configuration: %+v", *migration.Configuration())
	mig, err := migrator.New(migration.Configuration())
	if err != nil {
		writeMeta(err)
		log.Fatalf("failed to start Erda MySQL Migration: %v", err)
	}
	if err = mig.Run(); err != nil {
		writeMeta(err)
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Infoln("migrate complete !")
	writeMeta(err)

	os.Exit(0)
}

func startSandbox() {
	log.Infoln("create sandbox")
	sandbox := exec.Command("/usr/bin/run-mysqld")
	if err := sandbox.Start(); err != nil {
		log.Fatalf("failed to Start sandbox, err: %v", err)
	}
	if err := sandbox.Wait(); err != nil {
		log.Fatalf("failed to exec /usr/bin/run-mysqld, err: %v", err)
	}
}

func writeMeta(err error) {
	var data = `{"metadata": [{"name": "success", "value": "%s"}, {"name": "err", "value": "%s"}]}`
	if err == nil {
		data = fmt.Sprintf(data, "true", "nil")
	} else {
		data = fmt.Sprintf(data, "false", err.Error())
	}
	filename := migration.Configuration().MetaFilename()
	_ = ioutil.WriteFile(filename, []byte(data), 0644)
}
