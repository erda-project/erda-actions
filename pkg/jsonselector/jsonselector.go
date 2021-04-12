package jsonselector

import (
	"fmt"
	"strconv"

	"github.com/bitly/go-simplejson"
)

type Token struct {
	v string
	f string
}

func Select(data []byte, exp string) (string, error) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		return "", err
	}
	ets, err := split(exp)
	if err != nil {
		return "", err
	}
	return doSelect(js, ets)
}

func doSelect(js *simplejson.Json, ets []*Token) (string, error) {
	for _, t := range ets {
		switch t.f {
		case ".":
			js = js.Get(t.v)
		case "[":
			idx, err := strconv.Atoi(t.v)
			if err != nil {
				return "", err
			}
			js = js.GetIndex(idx)
		default:
			return "", fmt.Errorf("not supported f: %s", t.f)
		}
	}
	return fmt.Sprintf("%v", js.Interface()), nil
}

func split(str string) ([]*Token, error) {
	if str == "" {
		return nil, nil
	}
	if str[0] != '.' && str[0] != '[' {
		str = "." + str
	}
	l := len(str)
	vt := make([]*Token, 0)
	for i := 0; i < l; i++ {
		switch str[i] {
		case '.':
			vt = append(vt, &Token{f: "."})
		case '[':
			vt = append(vt, &Token{f: "["})
		case ']':
		default:
			t := vt[len(vt)-1]
			t.v = t.v + string(str[i])
		}
	}
	return vt, nil
}
