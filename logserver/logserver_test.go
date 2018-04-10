package main

import (
	"fmt"
	"testing"
	"io"
	"bytes"
)

var out io.Writer
func TestPrintLogs(t *testing.T) {

	var tests = []struct {
		logs []LogEntry
		want  string
	}{
		{nil, msgdatefmt},
	}

	for _, test := range tests {
		out = new(bytes.Buffer)
		descr := fmt.Sprintf("printLogs(out, %v)", test.logs)

		printLogs(out, test.logs)
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}
