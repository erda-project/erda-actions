package meta

import (
	"fmt"
	"os"

	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/strutil"
)

type ResultMeta struct {
	Key   string
	Value string
}

type Collector interface {
	Collect(meta ResultMeta)
	Store() error
}

type logCollector struct{}

func (l *logCollector) Collect(meta ResultMeta) {
	fmt.Fprintf(os.Stdout, "action meta: %s=%s\n", meta.Key, meta.Value)
}

func (l *logCollector) Store() error {
	return nil
}

type fileCollector struct {
	metaFilePath string
	metas        []ResultMeta
}

func (f *fileCollector) Collect(meta ResultMeta) {
	f.metas = append(f.metas, meta)
}

func (f *fileCollector) Store() error {
	if f.metaFilePath == "" {
		return nil
	}

	if len(f.metas) == 0 {
		return nil
	}
	var kvs []string
	for _, meta := range f.metas {
		kvs = append(kvs, fmt.Sprintf("%s=%s", meta.Key, meta.Value))
	}
	content := strutil.Join(kvs, "\n", true)
	return filehelper.CreateFile(f.metaFilePath, content, 0644)
}

type ResultMetaNotifier interface {
	Add(meta ResultMeta)
	Store() error
	register(c Collector)
}

type ResultMetaCollector struct {
	collectors []Collector
}

func (r *ResultMetaCollector) Add(key, value string) {
	for _, c := range r.collectors {
		c.Collect(ResultMeta{key, value})
	}
}

func (r *ResultMetaCollector) Store() error {
	for _, c := range r.collectors {
		if err := c.Store(); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultMetaCollector) Register(c Collector) {
	r.collectors = append(r.collectors, c)
}

type Option func(*ResultMetaCollector)

func WithFileCollector(metaFilePath string) Option {
	return func(c *ResultMetaCollector) {
		c.Register(&fileCollector{
			metaFilePath: metaFilePath,
			metas:        make([]ResultMeta, 0),
		})
	}
}

// NewResultMetaCollector creates a new ResultMetaCollector.
// default use log collector
func NewResultMetaCollector(opts ...Option) *ResultMetaCollector {
	c := &ResultMetaCollector{
		collectors: make([]Collector, 0),
	}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.collectors) == 0 {
		c.Register(&logCollector{})
	}
	return c
}
