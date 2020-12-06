// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
)

// Flags contains the command line parameters.
type Flags struct {
	Output    string
	JSONLines bool
	Template  string
	Args      []string
}

// Parse parses the command line parameters.
func (f *Flags) Parse(args []string) error {
	var unsafe bool

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&f.Output, "o", "-", "output `file`, defaults to stdout")
	flags.BoolVar(&f.JSONLines, "jsonl", false, "read JSON-line encoded files")
	flags.BoolVar(&unsafe, "unsafe", false, "enable unsafe version of the execute function")
	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "Usage: %s [options] mdtemplate [ jsonFile* ]\n", args[0])
		flag.PrintDefaults()
	}
	flags.Parse(args[1:])

	// enable unsafe version of execute function.
	if unsafe {
		templateFuncs["execute"] = executeUnsafe
	}

	// capture the remaining arguments
	args = flags.Args()
	if len(args) == 0 {
		return errors.New("error: missing template file")
	}

	f.Template = args[0]
	f.Args = args[1:]
	return nil
}
