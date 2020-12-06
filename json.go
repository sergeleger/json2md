// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io"
)

// JSON structure has file metadata and the JSON content of the file. It is exposed to the
// template execution.
type JSON struct {
	Filename string
	Content  interface{}
}

// ReadFrom reads single JSON object from reader
func (j *JSON) ReadFrom(r io.Reader) error {
	return json.NewDecoder(r).Decode(&j.Content)
}

// NewJSONReader creates a specialized JSON reader. The option jsonL indicates that
// encoding contains one JSON object per line.
func NewJSONReader(dest *JSON, jsonL bool) ReaderFrom {
	if !jsonL {
		return dest
	}

	return jsonLines{dest}
}

// jsonLines reads a JSON Lines formated file.
type jsonLines struct {
	dest *JSON
}

// read multiple json objects from the same file
func (j jsonLines) ReadFrom(r io.Reader) error {
	var data []interface{}
	var err error

	dec := json.NewDecoder(r)
	for {
		var d interface{}
		err = dec.Decode(&d)
		if err != nil {
			break
		}

		data = append(data, d)
	}

	if err == io.EOF {
		err = nil
	}

	j.dest.Content = data
	return err
}

// ReaderFrom defines an object that can reads its content from a io.Reader
type ReaderFrom interface {
	ReadFrom(r io.Reader) error
}
