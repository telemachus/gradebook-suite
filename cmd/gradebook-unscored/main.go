// Gb provides commands to work with student grades.
package main

import (
	"os"

	"github.com/telemachus/gradebook-suite/internal/cli"
)

func main() {
	os.Exit(cli.GradebookUnscored(os.Args[1:]))
}
