package cli

import (
	"fmt"

	"github.com/telemachus/gradebook"
)

// GradebookUnscored counts and prints unscored assignments for a class.
func GradebookUnscored(args []string) int {
	cmd := cmdFrom("gradebook-unscored", unscoredUsage)

	return runCommand(cmd, args, commandRun[string]{
		parse:     (*cmdEnv).parseUnscored,
		loadClass: true,
		action: func(cmd *cmdEnv, class *gradebook.Class, term string) {
			cmd.findTerm(class, term)
			cmd.loadUnscored(class, term)
			cmd.printUnscored(class)
		},
	})
}

func (cmd *cmdEnv) parseUnscored(args []string) string {
	og := cmd.commonOptsGroup(parseOpts{})

	term := ""
	og.String(&term, "term", "")

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
		fmt.Fprintln(cmd.stderr, cmd.usage)

		return ""
	}

	return term
}

func (cmd *cmdEnv) loadUnscored(class *gradebook.Class, term string) {
	if cmd.noOp() {
		return
	}

	err := class.LoadUnscored(cmd.directory, class.TermsByID[term])
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
	}
}

func (cmd *cmdEnv) printUnscored(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	students := class.StudentsSortedByName()
	for _, s := range students {
		fmt.Fprintf(cmd.stdout, "%s %s:\n", s.FirstName, s.LastName)

		for _, cat := range class.AssignmentCategoriesSortedByLabel() {
			n := s.UnscoredByCategory[cat]
			label := class.LabelsByAssignmentCategory[cat]
			word := "assignments"
			if n == 1 {
				word = "assignment"
			}
			fmt.Fprintf(cmd.stdout, "\t%s: %d unscored %s\n", label, n, word)
		}
	}
}
