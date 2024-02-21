package main

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

type testEnv struct {
	args           string
	stdout, stderr bytes.Buffer
}

func (e *testEnv) run() error {
	s := flag.NewFlagSet("monitor", flag.ContinueOnError)
	s.SetOutput(&e.stderr)
	return run(s, strings.Fields(e.args), &e.stdout)
}

func TestRun(t *testing.T) {
	happy := map[string]struct{ in, out string }{
		"normal": {
			"--filepath test.log --slack --storefile ~/Documents/index.json"
		}
	}

}
