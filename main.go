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
	"text/template"
)

var templateFuncs = template.FuncMap{}

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainErr() error {
	output := flag.String("o", "-", "output `file`, defaults to stdout")
	multiline := flag.Bool("jsonml", false, "attempts to read all json files as multi-line files")
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
		data, err = readJSON(args[0], *multiline)
		if err != nil {
			return err
		}

	case n > 1:
		obj := make([]interface{}, n)
		for i := range args {
			obj[i], err = readJSON(args[i], *multiline)
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
func readJSON(file string, multiline bool) (interface{}, error) {
	f, err := openFile(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read single valid json object from the file
	if !multiline {
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
