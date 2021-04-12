package parser

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/tap"
)

func ParserTapSuite(content string) (*tap.Testsuite, error) {
	var (
		p     *tap.Parser
		suite *tap.Testsuite
		err   error
	)
	r := strings.NewReader(content)

	if p, err = tap.NewParser(r); err != nil {
		return nil, errors.Wrapf(err, "create a new tap parser error")
	}

	if suite, err = p.Suite(); err != nil {
		return nil, errors.Wrapf(err, "get test suites error")
	}

	return suite, nil
}
