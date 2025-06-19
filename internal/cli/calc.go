package cli

import (
	"fmt"

	"github.com/telemachus/gradebook"
)

// GradebookCalc calculates and prints the grades for a class.
func GradebookCalc(args []string) int {
	cmd := cmdFrom("gradebook-calculate", calcUsage, suiteVersion)

	extraArgs, term := cmd.parseCalculate(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.findTerm(class, term)
	cmd.loadGrades(class, term)
	cmd.printAll(class)

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

func (cmd *cmdEnv) findTerm(class *gradebook.Class, term string) {
	if cmd.noOp() || term == "" {
		return
	}

	if _, ok := class.TermsByID[term]; !ok {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %q is not a valid term\n", cmd.name, term)
	}
}

func (cmd *cmdEnv) loadGrades(class *gradebook.Class, term string) {
	if cmd.noOp() {
		return
	}

	err := class.LoadGrades(cmd.directory, class.TermsByID[term])
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
	}
}

func (cmd *cmdEnv) printAll(class *gradebook.Class) {
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
