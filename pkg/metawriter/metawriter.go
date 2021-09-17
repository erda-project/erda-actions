// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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

// WriteKV is the shortcut for New(filename).Write(key, value)
func WriteKV(k string, v interface{}) error {
	return New(os.Getenv(metafile)).WriteKV(k, v)
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

func (w Writer) WriteKV(k string, v interface{}) error {
	data, err := json.Marshal(meta{Metadata: []ele{{Name: k, Value: fmt.Sprintf("%v", v)}}})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.filename, data, 0644)
}
