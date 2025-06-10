package cli

var calculateUsage = "usage: gradebook-calculate: TODO"

// GradebookCalculate calculates and displays the grades for a class.
func GradebookCalculate(_ []string) int {
	cmd := cmdFrom("gradebook-calculate", calculateUsage, suiteVersion)

	return cmd.exitValue
}
