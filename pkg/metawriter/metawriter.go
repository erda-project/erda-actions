package metawriter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const metafile = "METAFILE"

type meta struct {
	Metadata []ele `json:"metadata"`
}

type ele struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Write is the shortcut for New(filename).Write(key, value)
func Write(key string, value interface{}) (err error) {
	return New(os.Getenv(metafile)).Write(key, value)
}

type Writer struct {
	filename string
}

func New(filename string) *Writer {
	return &Writer{filename: filename}
}

func (w Writer) Write(key string, value interface{}) error {
	data, err := json.Marshal(meta{Metadata: []ele{{Name: key, Value: fmt.Sprintf("%v", value)}}})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.filename, data, 0644)
}
