package configs

import (
	"fmt"
	"testing"
	"io"
	"bytes"
	"strings"
)

var out io.Writer

func TestReadConfFile(t *testing.T) {

	var outconf = Conf{}
	
	// can easily add test cases - be careful about brittle tests
	var tests = []struct {
		filename string
		c *Conf
		want  string
	}{
		{"../logserver/conf.json", &outconf, "mainapp3.log"},
		{"", nil, ""},
		{"/path/to/no/file", nil, ""},
	}

	for _, test := range tests {
		out = new(bytes.Buffer)
		descr := fmt.Sprintf("ReadConfFile(%s, %v)", test.filename, test.c)

		err := ReadConfFile(test.filename, test.c)
		if err != nil {
			if test.want != "" {
				t.Errorf("%s failed: %v", descr, err)
			}
			continue
		}

		fmt.Fprint(out, outconf)
		got := out.(*bytes.Buffer).String()
		if !strings.Contains(got, test.want) {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}
