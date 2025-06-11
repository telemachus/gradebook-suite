package cli

var newUsage = "usage: gradebook-new: TODO"

// GradebookNew creates a new gradebook file for a class.
func GradebookNew(args []string) int {
	cmd := cmdFrom("gradebook-new", newUsage, suiteVersion)

	extraArgs, gbNew := cmd.parseNew(args)
	cmd.check(extraArgs)
	cmd.printHelpOrVersion()

	cmd.resolvePaths()
	class := cmd.unmarshalClass()
	cmd.checkNew(gbNew, class)
	// cmd.newGradebook(class, gbNew)

	return cmd.exitValue
}
