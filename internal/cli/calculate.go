package cli

import (
	"fmt"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
	"github.com/telemachus/gradebook-suite/internal/opts"
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
	cmd.calculateAndDisplay(term, class)

	return cmd.exitValue
}

func (cmd *cmdEnv) parseCalculate(args []string) ([]string, string) {
	og := opts.NewGroup(cmd.name)
	og.String(&cmd.classFile, "class", "class.json")
	og.String(&cmd.directory, "directory", "")
	og.Bool(&cmd.helpWanted, "help")
	og.Bool(&cmd.helpWanted, "h")
	og.Bool(&cmd.versionWanted, "version")

	term := ""
	og.String(&term, "term", "")

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return nil, ""
	}

	return og.Args(), term
}

func (cmd *cmdEnv) calculateAndDisplay(term string, class *gradebook.Class) {
	if cmd.noOp() {
		return
	}

	if err := class.LoadGrades(cmd.directory, class.TermsByID[term]); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)

		return
	}

	students := class.StudentsSortedByName()
	for _, s := range students {
		fmt.Fprintf(cmd.stdout, "%s %s\n", s.FirstName, s.LastName)
		totalAvg, err := s.TotalAverage(class.WeightsByAssignmentCategory)
		if err != nil {
			fmt.Fprintf(cmd.stderr, "\t%s: problem calculating total average: %s\n", cmd.name, err)

			return
		}

		fmt.Fprintf(cmd.stdout, "\tOverall average: %s\n", totalAvg)

		for _, cat := range class.AssignmentCategoriesSortedByLabel() {
			label := class.LabelsByAssignmentCategory[cat]

			avg, err := s.Average(cat)
			if err != nil {
				fmt.Fprintf(cmd.stderr, "\t%s: problem calculating %s: %s\n", cmd.name, label, err)

				return
			}

			fmt.Fprintf(cmd.stdout, "\t%s: %s\n", label, avg)

		}
	}
}
