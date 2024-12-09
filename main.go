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
	runGoInit := os.Getenv("GEN_INIT") == "1" || os.Getenv("GEN_INIT") == "true"
	flag.StringVar(&rflg.Command, "cmd", os.Getenv("GEN_CMD"), "executable `name`, i.e. 'foo', if you will run it as './foo help'")
	flag.StringVar(&rflg.CommandVariable, "var", os.Getenv("GEN_VAR"), "main command variable `name`, i.e. 'FooCommand', must be exported")
	flag.StringVar(&rflg.Package, "pkg", os.Getenv("GEN_PKG"), "Command package name, i.e. 'github.com/you/yourpackage/cmd/slackdump'")
	flag.BoolVar(&outflg.RunGoInit, "init", runGoInit, "run 'go init' in the output directory with the package name")
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
	RunGoInit bool
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
		if r.CommandVariable[0] < 'a' || r.CommandVariable[0] > 'z' {
			r.CommandVariable = strings.ToTitle(r.CommandVariable)
		} else {
			return errors.New("command variable must start with an uppercase letter")
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
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export"))
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
