package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var _ = loadDotEnv()

var (
	rflg   renderer
	outflg output
)

func init() {
	flag.StringVar(&rflg.Command, "cmd", os.Getenv("GEN_CMD"), "executable `name`, i.e. 'foo', if you will run it as './foo help'")
	flag.StringVar(&rflg.CommandVariable, "var", os.Getenv("GEN_VAR"), "main command variable `name`, i.e. 'FooCommand', must be exported")
	flag.StringVar(&rflg.Package, "pkg", os.Getenv("GEN_PKG"), "Command package name, i.e. 'github.com/you/yourpackage/cmd/slackdump'")
}

func main() {
	flag.Parse()

	outflg.OutputDir = os.Getenv("GEN_OUTPUT_DIR")
	if flag.NArg() == 0 && outflg.OutputDir == "" {
		flag.Usage()
		log.Fatal("invalid parameters")
	}
	if outflg.OutputDir == "" {
		outflg.OutputDir = flag.Arg(0)
	}

	if err := generate(rflg, outflg); err != nil {
		log.Fatal(err)
	}
}

type output struct {
	OutputDir string
}

func (o *output) validate() error {
	if o.OutputDir == "" {
		return errors.New("output directory cannot be empty")
	}
	return nil
}

// renderer defines replacements.
type renderer struct {
	Command         string
	CommandVariable string
	Package         string
}

func (r *renderer) validate() error {
	if r.Command == "" {
		return errors.New("command cannot be empty")
	}
	if r.CommandVariable == "" {
		return errors.New("command variable cannot be empty")
	}
	if r.CommandVariable[0] < 'A' || r.CommandVariable[0] > 'Z' {
		if 'a' <= r.CommandVariable[0] && r.CommandVariable[0] <= 'z' {
			r.CommandVariable = strings.ToUpper(r.CommandVariable[:1]) + r.CommandVariable[1:]
		} else {
			return errors.New("command variable must start with a letter")
		}
	}
	if r.Package == "" {
		return errors.New("package cannot be empty")
	}
	if r.Package[len(r.Package)-1] != '/' {
		r.Package += "/"
	}
	return nil
}

func (r *renderer) replacer() *strings.Replacer {
	return strings.NewReplacer("$$Command$$", r.Command, "MAIN__COMMAND", r.CommandVariable, "github.com/rusq/cmdgen/template/", r.Package)
}

const maxlines = 1 << 10

func loadDotEnv() error {
	f, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	lineno := 0
	for s.Scan() && lineno < maxlines {
		lineno++
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if after, ok := strings.CutPrefix(line, "export"); ok {
			line = strings.TrimSpace(after)
		}
		name, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("line %d: expected key=value pair", lineno)
		}
		if err := os.Setenv(name, value); err != nil {
			return err
		}
	}
	return nil
}
