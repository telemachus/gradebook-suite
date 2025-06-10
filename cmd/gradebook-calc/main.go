// Gb provides commands to work with student grades.
package main

import (
	"os"

	"github.com/telemachus/gradebook-suite/internal/cli"
)

func main() {
	os.Exit(cli.GradebookCalculate(os.Args[1:]))
}
