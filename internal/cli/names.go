package cli

import (
	"fmt"

	"github.com/telemachus/gradebook"
)

// GradebookNames prints the names of students in a class. The default
// output is "FirstName LastName", but the user can opt for "LastName,
// FirstName" instead.
func GradebookNames(args []string) int {
	cmd := cmdFrom("gradebook-names", namesUsage)

	return runCommand(cmd, args, commandRun[noArgs]{
		parse:     (*cmdEnv).parseNames,
		loadClass: true,
		action: func(cmd *cmdEnv, class *gradebook.Class, _ noArgs) {
			cmd.printNames(class)
		},
	})
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
