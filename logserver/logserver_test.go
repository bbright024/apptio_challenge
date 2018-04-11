package main

import (
	"fmt"
	"testing"
	"io"
	"bytes"
	"os"
)

var out io.Writer
func TestPrintLogs(t *testing.T) {

	// can easily add test cases - be careful about brittle tests
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

func TestConvertLogFile(t *testing.T) {

	var tests = []struct {
		file *os.File
		want []LogEntry
	}{
		{nil, nil},
	}

	for _, test := range tests {

		descr := fmt.Sprintf("convertLogFile(%v)", test.file)
		lentries := convertLogFile(test.file)		

		
		if lentries == nil && test.want == nil {
			continue
		} else if len(lentries) != len(test.want){
			t.Errorf("%s = %q, want %q", descr, lentries, test.want)
		} else {
			for i, le := range lentries {
				testle := test.want[i]
				// these tests will help check parsing 
				if le.Logtime != testle.Logtime {
					t.Errorf("%s = %q, want %q", descr, le.Logtime, testle.Logtime)
				}

				if le.Message != testle.Message {
					t.Errorf("%s = %q, want %q", descr, le.Message, testle.Message)
				}
			}				
		}
	}
}
