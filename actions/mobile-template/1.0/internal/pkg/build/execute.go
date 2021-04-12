package build

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/erda-project/erda-actions/actions/mobile-template/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
)

// Execute Generate mobile app code template, including ios, android, h5 & weixin mini program
func Execute() error {
	// trnw-cli init <projectName> --displayName demo --bundleId io.terminus.demo --packageName io.terminus.demo
	var cfg conf.Conf
	envconf.MustLoad(&cfg)
	fmt.Fprintln(os.Stdout, "sucessfully loaded action config")

	if cfg.ProjectName == "" {
		cfg.ProjectName = "templates"
	}

	templteGenCmd := exec.Command("trnw-cli", "init", cfg.ProjectName,
		"--type", "default",
		"--displayName", cfg.DisplayName,
		"--iosBundleId", cfg.BundleID,
		"--androidPackageName", cfg.PackageName)
	fmt.Println(templteGenCmd.Args)
	templteGenCmd.Stdout = os.Stdout
	templteGenCmd.Stderr = os.Stderr
	if err := templteGenCmd.Run(); err != nil {
		return err
	}

	return nil
}
