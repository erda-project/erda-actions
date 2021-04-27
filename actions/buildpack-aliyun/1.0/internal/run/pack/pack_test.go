package pack

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParseMemory(t *testing.T) {

	memory := strconv.FormatFloat(float64(2048.1024*1000000), 'f', 0, 64)

	fmt.Println(memory)

}
