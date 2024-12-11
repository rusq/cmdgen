# CMDGEN

An utility to generate command line interface for Go programs.

The template uses the stripped down version of the command processing
logic taken from Go command.  It is intended to be a lightweight
alternative to urfave/cli and cobra.

## Why?

Why would you want to use this instead of the aforementioned
libraries?

1. Easily create new subcommands to the main command with unlimited
   depth.  For example, if your command is `foo`, it can have
   subcommands `foo bar`, `foo bar baz`, etc.
2. Centralised configuration storage (package cfg) that can be
   partitioned by flag mask (see template/internal/cfg/cfg.go).
3. Extensibility.  As generated files are part of your package you can
   easily modify them to suit your needs (for example, see
   [slackdump]).
4. Zero dependencies.
5. Lightweight and simple.

Why would you not want to use this package?

1. You don't like `flag` package
2. You fancy those double dash `--flag` arguments.
3. You don't want to bother with extending it.
4. You need something more complex, like binding the environment
   variables to struct fields.

[slackdump]: https://github.com/rusq/slackdump/blob/master/cmd/slackdump

## Usage

Install the command:
```shell
go install github.com/rusq/cmdgen
```

Run it:
```shell
cmdgen -cmd yourmaincmd -var YourMainCmd -pkg github.com/you/yourpackage/cmd/yourmaincmd /path/to/your/project/cmd/yourmaincmd
```

## Command line switches

```
Usage: cmdgen <flags> <path>
  -cmd name
        executable name, i.e. 'foo', if you will run it as './foo help'
  -pkg string
        Command package name, i.e. 'github.com/you/yourpackage/cmd/slackdump'
  -var name
        main command variable name, i.e. 'FooCommand', must be exported
```

## Environment variables

```
GEN_CMD=foo
GEN_PKG=github.com/you/foo
GEN_VAR=FooCommand
GEN_OUTPUT_DIR=/path/to/your/project/cmd/foo
```

# Tutors

Conventions:
- For simplicity, I'm going to call the package that we're extending "foobar",
  i.e. "github.com/acme/foobar/cmd/foobar".
- Command name is "foobar"
- Main command variable is "Foobar"

## Adding new command

We are going to add a "whambam" subcommand to the "foobar" command, so
that user can run it as `foobar whambam`, and when help is printed,
user sees:

```
Usage: foobar subcommand [flags]

    whambam  - book some time at WHAMBAM HOTEL

```


1. Create a new directory under the cmd/foobar/internal directory, i.e.
   cmd/foobar/internal/whambam.
2. Create a file named "whambam.go" with the following contents:
   ```go
   package whambam
   
   import "github.com/acme/foobar/cmd/foobar/internal/golang/base"
   
   var CommandWhamBam = &base.Command{
	   UsageLine:  "foobar whambam [flags]",
	   Short:      "sends a web request to book time at Whambam Hotel",
	   PrintFlags: true,
	   Run:        runWhambam,
   }
   
   func runWhambam(ctx context.Context, cmd *base.Command, args []string) error {
	   return errors.New("implement me")
   }
   ```
3. Add it to the slice of commands at "main.go:21":
   ```go
   //...
   func init() {
	   base.Foobar.Commands = []*base.Command{
		   whambam.CommandWhamBam,
	   }
   }
   //...
   ```
4. Save all files, if you haven't done so.
5. Run `go run ./cmd/foobar help` and observe that "whambam" is now
   in the commands list.
6. Run `go run ./cmd/foobar whambam` to see the error that we
   carefully planted there.
