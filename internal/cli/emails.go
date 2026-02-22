package cli

import (
	"fmt"
	"strings"

	"github.com/telemachus/gradebook"
)

// GradebookEmails prints the emails of students in a class.
func GradebookEmails(args []string) int {
	cmd := cmdFrom("gradebook-emails", emailsUsage)

	return runCommand(cmd, args, commandRun[noArgs]{
		parse:     (*cmdEnv).parseNoArgs,
		loadClass: true,
		action: func(cmd *cmdEnv, class *gradebook.Class, _ noArgs) {
			cmd.printEmails(class)
		},
	})
}

func (cmd *cmdEnv) printEmails(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	emails := class.EmailsSortedByStudentName()
	fmt.Fprintln(cmd.stdout, strings.Join(emails, "\n"))
}
