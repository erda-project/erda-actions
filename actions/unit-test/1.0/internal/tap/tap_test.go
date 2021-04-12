package tap

import (
	"strings"
	"testing"

	"github.com/robertkrimen/terst"

	"github.com/erda-project/erda-actions/actions/unit-test/1.0/internal/tap"
)

func TestBasic(t *testing.T) {
	terst.Terst(t)

	r := strings.NewReader(`TAP version 13
1..2
ok 1
not ok 2`)
	p, e := tap.NewParser(r)
	terst.Is(e, nil, "No error parsing preamble")

	s, e := p.Suite()
	terst.Is(e, nil, "No error parsing input")

	terst.Is(len(s.Tests), 2, "Right number of tests")

	terst.Is(s.Tests[0].Ok, true, "First test ok")
	terst.Is(s.Tests[1].Ok, false, "Second test not ok")
}
