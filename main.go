// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"stringsReplace": strings.Replace,

	// sanitize table cell strings to remove newlines
	"tablecell": func(s string) string {
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\n", "<br>")
		return s
	},

	// concatenate joins the list of elements
	"concatenate": func(sep string, all ...interface{}) string {
		var sb strings.Builder

		n := 0
		for i := range all {
			if all[i] == nil {
				continue
			}

			n++
			if n > 1 {
				fmt.Fprint(&sb, sep)
			}

			fmt.Fprintf(&sb, "%v", all[i])
		}

		return sb.String()
	},
}

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainErr() error {
	output := flag.String("o", "-", "output `file`, defaults to stdout")
	jsonl := flag.Bool("jsonl", false, "read JSON-line encoded files")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] mdtemplate [ jsonFile* ]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		return errors.New("error: missing template file")
	}

	// read template source and create the template
	tmplFile := flag.Arg(0)
	content, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return err
	}

	tmpl := template.New(tmplFile)
	tmpl.Funcs(templateFuncs)
	_, err = tmpl.Parse(string(content))
	if err != nil {
		return err
	}

	// Read JSON source(s) data in memory
	n := flag.NArg() - 1
	args := flag.Args()[1:]
	var data interface{}
	switch {
	case n == 1:
		data, err = readJSON(args[0], *jsonl)
		if err != nil {
			return err
		}

	case n > 1:
		obj := make([]interface{}, n)
		for i := range args {
			obj[i], err = readJSON(args[i], *jsonl)
			if err != nil {
				return err
			}
		}
		data = obj
	}

	// Buffer the template's output
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}

	// Copy results to stdout
	w, err := createFile(*output)
	_, err = io.Copy(w, &buf)
	w.Close()
	return err
}

// readJSON reads a single file JSON file.
func readJSON(file string, jsonl bool) (interface{}, error) {
	f, err := openFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read single valid json object from the file
	if !jsonl {
		var data interface{}
		err = json.NewDecoder(f).Decode(&data)
		return data, err
	}

	// read multiple json objects from the same file
	var data []interface{}
	dec := json.NewDecoder(f)
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

	return data, err
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
