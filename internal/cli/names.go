package cli

import (
	"cmp"
	"fmt"
	"slices"

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

	students := make([]*gradebook.Student, 0, len(class.StudentsByEmail))
	for _, student := range class.StudentsByEmail {
		students = append(students, student)
	}
	slices.SortFunc(students, cmpStudent)

	for _, s := range students {
		switch cmd.lastFirst {
		case true:
			fmt.Fprintf(cmd.stdout, "%s, %s\n", s.LastName, s.FirstName)
		default:
			fmt.Fprintf(cmd.stdout, "%s %s\n", s.FirstName, s.LastName)
		}
	}
}

func cmpStudent(x, y *gradebook.Student) int {
	if x.LastName == y.LastName {
		return cmp.Compare(x.FirstName, y.FirstName)
	}

	return cmp.Compare(x.LastName, y.LastName)
}
