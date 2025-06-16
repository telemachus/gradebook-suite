package cli

import (
	"fmt"
	"strings"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

// GradebookEmails prints the emails of students in a class.
func GradebookEmails(args []string) int {
	cmd := cmdFrom("gradebook-emails", emailsUsage, suiteVersion)

	extraArgs := cmd.parse(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.printEmails(class)

	return cmd.exitValue
}

func (cmd *cmdEnv) printEmails(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	emails := class.EmailsSortedByStudentName()
	fmt.Fprintln(cmd.stdout, strings.Join(emails, "\n"))
}
