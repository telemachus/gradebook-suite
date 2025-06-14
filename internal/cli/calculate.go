package cli

import (
	"fmt"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

var calculateUsage = "usage: gradebook-calculate: TODO"

// GradebookCalculate calculates and displays the grades for a class.
func GradebookCalculate(args []string) int {
	cmd := cmdFrom("gradebook-calculate", calculateUsage, suiteVersion)

	extraArgs, term := cmd.parseCalculate(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.loadGrades(term, class)
	cmd.displayAll(class)

	return cmd.exitValue
}

func (cmd *cmdEnv) parseCalculate(args []string) ([]string, string) {
	og := cmd.commonOptsGroup()

	term := ""
	og.String(&term, "term", "")

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return nil, ""
	}

	return og.Args(), term
}

func (cmd *cmdEnv) loadGrades(term string, class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	err := class.LoadGrades(cmd.directory, class.TermsByID[term])
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
	}
}

func (cmd *cmdEnv) displayAll(class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	students := class.StudentsSortedByName()
	for _, s := range students {
		fmt.Fprintf(cmd.stdout, "%s %s\n", s.FirstName, s.LastName)

		fmt.Fprintf(cmd.stdout, "\tOverall average: %s\n", s.TotalAverage(class.WeightsByAssignmentCategory))

		for _, cat := range class.AssignmentCategoriesSortedByLabel() {
			fmt.Fprintf(cmd.stdout, "\t%s: %s\n", class.LabelsByAssignmentCategory[cat], s.Average(cat))
		}
	}
}
