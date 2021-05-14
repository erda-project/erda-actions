package main

type PackageManager interface {
	Analysis(cfg *Conf) (*ResultMetas, error)
}
