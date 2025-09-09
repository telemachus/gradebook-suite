package cli

import (
	"fmt"

	"github.com/telemachus/gradebook"
)

// GradebookNames prints the names of students in a class. The default
// output is "FirstName LastName", but the user can opt for "LastName,
// FirstName" instead.
func GradebookNames(args []string) int {
	cmd := cmdFrom("gradebook-names", namesUsage, suiteVersion)

	cmd.parse(args)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.printNames(class)

	return cmd.exitValue
}

func (cmd *cmdEnv) printNames(class *gradebook.Class) {
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
