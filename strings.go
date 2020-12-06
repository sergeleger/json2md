// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

// mdTableCell cleans the cell content to usable in a Markdown cell.
func mdTableCell(str string) string {
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "<br>")
	return str
}

// join joins a list of possibly nil elements.
func join(sep string, all ...interface{}) string {
	n := len(all) - 1

	var sb strings.Builder
	for i, v := range all {
		if v == nil {
			continue
		}

		fmt.Fprint(&sb, v)
		if i < n {
			sb.WriteString(sep)
		}
	}

	return sb.String()
}
