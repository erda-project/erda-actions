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
	"strconv"
)

const metafile = "METAFILE"

var (
	warnIndex uint64
	errIndex  uint64

	w = New(os.Getenv(metafile))
	m = meta{}
)

type meta struct {
	Metadata []ele `json:"metadata"`
}

type ele struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// Write is the shortcut for New(filename).Write(key, value)
func Write(values map[string]interface{}) (err error) {
	return w.Write(values)
}

// WriteKV is the shortcut for New(filename).Write(key, value)
func WriteKV(k string, v interface{}) error {
	return w.WriteKV(k, v)
}

// WriteSuccess writes the final result is whether success or fails.
func WriteSuccess(success bool) error {
	return w.WriteSuccess(success)
}

// WriteLink writes key-value with link
func WriteLink(k string, v interface{}) error {
	return w.WriteLink(k, v)
}

// WriteWarn writes warn info to meta file
func WriteWarn(v interface{}) error {
	return w.WriteWarn(v)
}

// WriteError writes err info to meta file
func WriteError(v interface{}) error {
	return w.WriteError(v)
}

type Writer struct {
	filename string
}

func New(filename string) *Writer {
	return &Writer{filename: filename}
}

func (w Writer) Write(values map[string]interface{}) error {
	for k, v := range values {
		m.Metadata = append(m.Metadata, ele{Name: k, Value: fmt.Sprintf("%v", v)})
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.filename, data, 0644)
}

func (w Writer) WriteKV(k string, v interface{}) error {
	return w.Write(map[string]interface{}{k: v})
}

// WriteSuccess writes the final result is whether success or fails.
func (w Writer) WriteSuccess(success bool) error {
	return w.WriteKV("success", success)
}

// WriteLink writes key-value with link
func (w Writer) WriteLink(k string, v interface{}) error {
	m.Metadata = append(m.Metadata, ele{Name: k, Value: fmt.Sprintf("%v", v), Type: "link"})
	return w.Write(make(map[string]interface{}))
}

// WriteWarn writes warn info to meta file
func (w Writer) WriteWarn(v interface{}) error {
	warnIndex++
	return w.WriteKV("warn-"+strconv.FormatUint(warnIndex, 10), v)
}

// WriteError writes err info to meta file
func (w Writer) WriteError(v interface{}) error {
	errIndex++
	return w.WriteKV("err-"+strconv.FormatUint(errIndex, 10), v)
}
