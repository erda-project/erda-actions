package confutil

import (
	"testing"

	"gotest.tools/assert"
)

type Address struct {
	Province string
	City     string
	Capital  bool
}

func TestStruct2Map(t *testing.T) {
	var st Address
	m := Struct2Map(st)
	assert.Equal(t, len(m), 3)
}
