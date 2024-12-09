# CMDGEN

An utility to generate command line interface for Go programs.

The template uses the stripped down version of the command processing logic
taken from Go command.  It is intended to be a lightweight alternative to
urfave/cli and cobra.

## Usage

Install the command:
```shell
go install github.com/rusq/cmdgen
```

Run it:
```shell
cmdgen -cmd yourmaincmd -pkg github.com/you/yourpackage/ -var YourMainCmd  /path/to/your/project/cmd/yourmaincmd
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
