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
func Write(m map[string]interface{}) (err error) {
	return New(os.Getenv(metafile)).Write(m)
}

type Writer struct {
	filename string
}

func New(filename string) *Writer {
	return &Writer{filename: filename}
}

func (w Writer) Write(m map[string]interface{}) error {
	var mt meta
	for k, v := range m {
		mt.Metadata = append(mt.Metadata, ele{k, fmt.Sprintf("%v", v)})
	}
	data, err := json.Marshal(mt)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.filename, data, 0644)
}
