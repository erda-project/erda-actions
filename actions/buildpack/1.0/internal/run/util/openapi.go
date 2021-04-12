package util

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

func GetDiceOpenAPIAddress() (string, error) {
	type envConfig struct {
		OpenAPIAddr string `env:"DICE_OPENAPI_ADDR,required"`
	}
	var envCfg envConfig
	err := env.Parse(&envCfg)
	if err != nil {
		return "", errors.Wrap(err, "unknown openapi address, cannot query or register artifact, skip")
	}
	return envCfg.OpenAPIAddr, nil
}
