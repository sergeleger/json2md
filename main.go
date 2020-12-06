// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	// string helper methods
	"stringsReplace": strings.Replace,
	"tablecell":      mdTableCell,
	"join":           join,

	// execute
	"execute": executeSafe,
}

func main() {
	if err := mainErr(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainErr(args []string) error {
	var flags Flags
	if err := flags.Parse(args); err != nil {
		return err
	}

	tmpl, err := readTemplate(flags.Template)
	if err != nil {
		return err
	}

	// Read JSON source(s) data in memory
	data := make([]JSON, len(flags.Args))
	for i, f := range flags.Args {
		data[i].Filename = f

		r, err := openFile(f)
		if err != nil {
			return err
		}

		err = NewJSONReader(&data[i], flags.JSONLines).ReadFrom(r)
		r.Close()
		if err != nil {
			return err
		}
	}

	// simplify the template object if we were called with a single input file.
	var tmplObj interface{} = data
	if len(data) == 1 {
		tmplObj = data[0]
	}

	// Buffer the template's output
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplObj)
	if err != nil {
		return err
	}

	// Copy results to stdout
	w, err := createFile(flags.Output)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, &buf)
	return w.Close()
}

// readTemplates reads and creates the template object.
func readTemplate(file string) (*template.Template, error) {
	r, err := openFile(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tmpl := template.New(file)
	tmpl.Funcs(templateFuncs)
	return tmpl.Parse(string(content))
}

// createFile opens the file for writing, if the file is "-" returns stdout.
func createFile(file string) (io.WriteCloser, error) {
	if file == "-" {
		return os.Stdout, nil
	}

	return os.Create(file)
}

// openFile opens a file for reading, if the file is "-" returns stdin.
func openFile(file string) (io.ReadCloser, error) {
	if file == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}

	return os.Open(file)
}
