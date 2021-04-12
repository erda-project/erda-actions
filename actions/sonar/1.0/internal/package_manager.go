package main

type PackageManager interface {
	Analysis(cfg *Conf) (map[ResultKey]string, error)
}
