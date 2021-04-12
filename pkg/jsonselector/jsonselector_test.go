package jsonselector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	var (
		v   string
		err error
	)
	v, err = Select([]byte(`{"data": {"id": 123}}`), "data.id")
	require.NoError(t, err)
	require.Equal(t, "123", v)

	v, err = Select([]byte(`
{
  "data": {
    "list": [
      {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
      {
        "apis": [
          {
            "apiInfo": "test it"
          }
        ]
      }
    ]
  }
}
`), "data.list[10].apis[0].apiInfo")
	require.NoError(t, err)
	require.Equal(t, "test it", v)
}

func Test_split(t *testing.T) {
	var (
		ts  []*Token
		err error
	)

	ts, err = split("data.id")
	require.NoError(t, err)
	require.Equal(t, []*Token{
		{f: ".", v: "data"},
		{f: ".", v: "id"},
	}, ts)

	ts, err = split("data.list[10].apis[0].apiInfo")
	require.NoError(t, err)
	require.Equal(t, []*Token{
		{f: ".", v: "data"},
		{f: ".", v: "list"},
		{f: "[", v: "10"},
		{f: ".", v: "apis"},
		{f: "[", v: "0"},
		{f: ".", v: "apiInfo"},
	}, ts)
}
