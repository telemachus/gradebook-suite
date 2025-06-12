package cli

import (
	"fmt"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

var namesUsage = "usage: gradebook-names: TODO"

// GradebookNames displays the names of students in a class. The default
// display is "FirstName LastName", but the user can opt for "LastName,
// FirstName" instead.
func GradebookNames(args []string) int {
	cmd := cmdFrom("gradebook-names", namesUsage, suiteVersion)

	extraArgs := cmd.parse(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.displayNames(class)

	return cmd.exitValue
}

func (cmd *cmdEnv) displayNames(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	students := class.StudentsSortedByName()
	for _, s := range students {
		switch cmd.lastFirst {
		case true:
			fmt.Fprintf(cmd.stdout, "%s, %s\n", s.LastName, s.FirstName)
		default:
			fmt.Fprintf(cmd.stdout, "%s %s\n", s.FirstName, s.LastName)
		}
	}
}
