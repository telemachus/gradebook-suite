package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

var emailsUsage = "usage: gradebook-emails: TODO"

// GradebookEmails displays the emails of students in a class.
func GradebookEmails(args []string) int {
	cmd := cmdFrom("gradebook-emails", emailsUsage, suiteVersion)

	extraArgs := cmd.parse(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.displayEmails(class)

	return cmd.exitValue
}

func (cmd *cmdEnv) displayEmails(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	emails := make([]string, 0, len(class.StudentsByEmail))
	for email := range class.StudentsByEmail {
		emails = append(emails, email)
	}

	slices.SortFunc(emails, func(emailA, emailB string) int {
		studentA := class.StudentsByEmail[emailA]
		studentB := class.StudentsByEmail[emailB]

		return cmpStudent(studentA, studentB)
	})

	fmt.Fprintln(cmd.stdout, strings.Join(emails, "\n"))
}
