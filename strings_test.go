// Copyright 2020 Serge Léger. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package main

import "testing"

func TestJoin(t *testing.T) {
	data := []interface{}{"a", nil, "b", 1}

	got := join(",", data...)
	expected := "a,b,1"
	if got != expected {
		t.Fatalf("expected: %q, got: %q", expected, got)
	}
}
