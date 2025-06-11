// Package cli creates and runs a command line interface.
package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
	"github.com/telemachus/gradebook-suite/internal/opts"
)

const (
	exitSuccess    = 0
	exitFailure    = 1
	suiteClassFile = "class.json"
	suiteDirectory = "."
	suiteVersion   = "v0.0.1"
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

func cmdFrom(name, usage, version string) *cmdEnv {
	return cmdFromWithWriters(name, usage, version, os.Stdout, os.Stderr)
}

func cmdFromWithWriters(name, usage, version string, stdout, stderr io.Writer) *cmdEnv {
	return &cmdEnv{
		exitValue: exitSuccess,
		name:      name,
		usage:     usage,
		version:   version,
		stdout:    stdout,
		stderr:    stderr,
	}
}

func (cmd *cmdEnv) parse(args []string) []string {
	og := opts.NewGroup(cmd.name)
	og.String(&cmd.classFile, "class", "class.json")
	og.String(&cmd.directory, "directory", "")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.versionWanted, "version")

	// TODO: add switch here for additional flags for gb-calc.
	switch cmd.name {
	case "gradebook-names":
		og.Bool(&cmd.lastFirst, "last-first")
	}

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return nil
	}

	return og.Args()
}

type newCfg struct {
	gbName string
	gbType string
	gbDate string
}

func (cmd *cmdEnv) parseNew(args []string) ([]string, newCfg) {
	og := opts.NewGroup(cmd.name)
	og.String(&cmd.classFile, "class", "class.json")
	og.String(&cmd.directory, "directory", "")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.versionWanted, "version")

	var gbCfg newCfg
	og.String(&gbCfg.gbName, "name", "")
	og.String(&gbCfg.gbType, "type", "")
	og.String(&gbCfg.gbDate, "date", "")

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return nil, gbCfg
	}

	return og.Args(), gbCfg
}

func (cmd *cmdEnv) check(extraArgs []string) {
	if cmd.minNoOp() {
		return
	}

	numExtraArgs := len(extraArgs)
	if numExtraArgs != 0 {
		cmd.exitValue = exitFailure

		var s string
		if numExtraArgs > 1 {
			s = "s"
		}

		fmt.Fprintf(cmd.stderr, "%s: unrecognized argument%s: %+v\n", cmd.name, s, extraArgs)
	}
}

func (cmd *cmdEnv) checkNew(cfg newCfg, class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	isValidName(cmd, cfg.gbName)
	isValidType(cmd, cfg.gbName, class)
	isValidDate(cmd, cfg.gbName)
}

func isValidName(cmd *cmdEnv, gbName string) {
	if gbName == "" || invalidGbNameRegex.MatchString(gbName) {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -name: %q\n", cmd.name, gbName)
	}
}

func isValidType(cmd *cmdEnv, gbType string, class *gradebook.Class) {
	if cmd.minNoOp() {
		return
	}

	if !slices.Contains(class.AssignmentTypes, gbType) {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -type: %q\n", cmd.name, gbType)
	}
}

func isValidDate(cmd *cmdEnv, gbDate string) {
	if cmd.minNoOp() {
		return
	}

	if _, err := time.Parse("20060102", gbDate); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -date: %q\n", cmd.name, gbDate)
	}
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
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

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
