// Package cli creates and runs a command line interface.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/telemachus/gradebook"
	"github.com/telemachus/opts"
)

const (
	exitSuccess    = 0
	exitFailure    = 1
	suiteClassFile = "class.json"
	suiteDirectory = "."
	suiteVersion   = "v20260222"
)

var invalidGbNameRegex = regexp.MustCompile(`[^A-Za-z0-9._-]`)

type cmdEnv struct {
	stdout        io.Writer
	stderr        io.Writer
	name          string
	classFile     string
	directory     string
	usage         string
	version       string
	exitValue     int
	lastFirst     bool
	helpWanted    bool
	versionWanted bool
}

type parseOpts struct {
	lastFirst bool
}

type noArgs struct{}

type commandRun[T any] struct {
	parse     func(*cmdEnv, []string) T
	action    func(*cmdEnv, *gradebook.Class, T)
	loadClass bool
}

func cmdFrom(name, usage string) *cmdEnv {
	return cmdFromWithWriters(name, usage, os.Stdout, os.Stderr)
}

func cmdFromWithWriters(name, usage string, stdout, stderr io.Writer) *cmdEnv {
	return &cmdEnv{
		exitValue: exitSuccess,
		name:      name,
		usage:     usage,
		version:   suiteVersion,
		stdout:    stdout,
		stderr:    stderr,
	}
}

func runCommand[T any](cmd *cmdEnv, args []string, runCfg commandRun[T]) int {
	parsed := runCfg.parse(cmd, args)
	cmd.printHelpOrVersion()
	if runCfg.loadClass {
		cmd.resolvePaths()
		class := cmd.unmarshalClass()
		runCfg.action(cmd, class, parsed)

		return cmd.exitValue
	}

	runCfg.action(cmd, nil, parsed)

	return cmd.exitValue
}

func (cmd *cmdEnv) parse(args []string) {
	cmd.parseWithOpts(args, parseOpts{})
}

func (cmd *cmdEnv) parseNames(args []string) noArgs {
	cmd.parseWithOpts(args, parseOpts{lastFirst: true})

	return noArgs{}
}

func (cmd *cmdEnv) parseNoArgs(args []string) noArgs {
	cmd.parse(args)

	return noArgs{}
}

func (cmd *cmdEnv) parseWithOpts(args []string, parseCfg parseOpts) {
	og := cmd.commonOptsGroup(parseCfg)

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
		fmt.Fprintln(cmd.stderr, cmd.usage)
	}
}

func (cmd *cmdEnv) commonOptsGroup(parseCfg parseOpts) *opts.Group {
	og := opts.NewGroup(cmd.name)
	og.String(&cmd.classFile, "class", "class.json")
	og.StringZero(&cmd.directory, "dir")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.versionWanted, "version")

	if parseCfg.lastFirst {
		og.Bool(&cmd.lastFirst, "last-first")
	}

	return og
}

func (cmd *cmdEnv) printHelpOrVersion() {
	if cmd.minNoOp() {
		return
	}

	switch {
	case cmd.helpWanted:
		fmt.Fprintln(cmd.stdout, cmd.usage)
	case cmd.versionWanted:
		fmt.Fprintf(cmd.stdout, "%s: %s\n", cmd.name, cmd.version)
	}
}

func (cmd *cmdEnv) resolvePaths() {
	if cmd.noOp() {
		return
	}

	absDirectory, err := filepath.Abs(cmd.directory)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return
	}

	absClassFile := filepath.Join(absDirectory, cmd.classFile)
	absClassFile, err = filepath.Abs(absClassFile)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return
	}

	cmd.directory = absDirectory
	cmd.classFile = absClassFile
}

func (cmd *cmdEnv) unmarshalClass() *gradebook.Class {
	if cmd.noOp() {
		return nil
	}

	class, err := gradebook.UnmarshalClass(cmd.classFile)
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: problem unmarshaling class: %s\n", cmd.name, err)

		return nil
	}
	if err = class.Validate(); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: problem validating class: %s\n", cmd.name, err)

		return nil
	}

	return class
}

func (cmd *cmdEnv) minNoOp() bool {
	return cmd.exitValue != exitSuccess
}

func (cmd *cmdEnv) noOp() bool {
	return cmd.exitValue != exitSuccess || cmd.helpWanted || cmd.versionWanted
}
