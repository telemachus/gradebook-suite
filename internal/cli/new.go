package cli

var newUsage = "usage: gradebook-new: TODO"

// GradebookNew creates a new gradebook file for a class.
func GradebookNew(_ []string) int {
	cmd := cmdFrom("gradebook-new", newUsage, suiteVersion)

	return cmd.exitValue
}
