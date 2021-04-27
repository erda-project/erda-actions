package conf

type conf struct {
	params       *params
	platformEnvs *platformEnv
	easyUse      *easyUse
}

var cfg conf

func Params() *params {
	return cfg.params
}
func PlatformEnvs() *platformEnv {
	return cfg.platformEnvs
}
func EasyUse() *easyUse {
	return cfg.easyUse
}

func Initialize() error {

	// init platform envs first, no depends
	_platformEnv, err := initPlatformEnvs()
	if err != nil {
		return err
	}
	cfg.platformEnvs = _platformEnv

	// init params depends on: platform envs
	_params, err := initParams()
	if err != nil {
		return err
	}
	cfg.params = _params

	// generate easy use last, depends: platform envs, params
	_easyUse, err := generateEasyUseLast()
	if err != nil {
		return err
	}
	cfg.easyUse = _easyUse

	return nil
}
