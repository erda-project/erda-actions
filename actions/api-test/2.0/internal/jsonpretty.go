package main

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

func jsonPretty(o interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("recover from jsonPretty: %v", r)
		}
	}()
	if o == nil {
		return ""
	}
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(o); err != nil {
		panic(err)
	}
	return buffer.String()
}

// jsonOneLine remove newline added by json encoder.Encode
func jsonOneLine(o interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("recover from jsonOneLine: %v", r)
		}
	}()
	if o == nil {
		return ""
	}
	switch o.(type) {
	case string: // 去除引号
		return o.(string)
	case []byte: // 去除引号
		return string(o.([]byte))
	default:
		var buffer bytes.Buffer
		enc := json.NewEncoder(&buffer)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(o); err != nil {
			panic(err)
		}
		return strings.TrimSuffix(buffer.String(), "\n")
	}
}
