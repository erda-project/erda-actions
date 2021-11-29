package pkg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
)

func migrationErda(cfg *conf.Conf) (string, error) {
	tmpTarPath := "/tmp/migration.tar.gz"

	tar := exec.Command("tar", "-zcf", tmpTarPath,
		"-C", cfg.MigrationDir, ".")
	tar.Stdout = os.Stdout
	tar.Stderr = os.Stderr
	if err := tar.Run(); err != nil {
		return "", err
	}

	r, err := UploadFileNew(tmpTarPath, *cfg)
	if err != nil {
		return "", err
	}
	if !r.Success {
		return "", fmt.Errorf(r.Error.Msg)
	}
	return r.Data.DownloadURL, nil
}
