// Package cfg contains common configuration variables.
package cfg

import (
	"flag"
	"log/slog"
	"os"
)

var (
	TraceFile   string
	LogFile     string
	JsonHandler bool
	Verbose     bool

	SomeFlag string

	Log *slog.Logger
)

type FlagMask uint16

const (
	DefaultFlags FlagMask = 0
	OmitSomeFlag FlagMask = 1 << (iota - 1)

	OmitAll = OmitSomeFlag
)

// SetBaseFlags sets base flags.
func SetBaseFlags(fs *flag.FlagSet, mask FlagMask) {
	fs.StringVar(&TraceFile, "trace", os.Getenv("TRACE_FILE"), "trace `filename`")
	fs.StringVar(&LogFile, "log", os.Getenv("LOG_FILE"), "log `file`, if not specified, messages are printed to STDERR")
	fs.BoolVar(&JsonHandler, "log-json", os.Getenv("JSON_LOG") != "", "log in JSON format")
	fs.BoolVar(&Verbose, "v", os.Getenv("DEBUG") != "", "verbose messages")

	if mask&OmitSomeFlag == 0 {
		fs.StringVar(&SomeFlag, "some-flag", "", "some flag sets something")
	}
}
