package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/telemachus/gradebook"
)

// GradebookNew creates a new gradebook file for a class.
func GradebookNew(args []string) int {
	cmd := cmdFrom("gradebook-new", newUsage)

	return runCommand(cmd, args, commandRun[newCfg]{
		parse:     (*cmdEnv).parseNew,
		loadClass: true,
		action: func(cmd *cmdEnv, class *gradebook.Class, gbCfg newCfg) {
			cmd.checkNew(class, gbCfg)
			cmd.newGradebook(class, gbCfg)
		},
	})
}

type newCfg struct {
	gbName string
	gbType string
	gbDate string
}

func (cmd *cmdEnv) parseNew(args []string) newCfg {
	og := cmd.commonOptsGroup(parseOpts{})

	var cfg newCfg
	og.String(&cfg.gbName, "name", "")
	og.String(&cfg.gbType, "type", "")
	og.String(&cfg.gbDate, "date", "")

	if err := og.Parse(args); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: %s\n", cmd.name, err)
		fmt.Fprintln(cmd.stderr, cmd.usage)

		return cfg
	}

	if cfg.gbDate == "" {
		now := time.Now()
		cfg.gbDate = now.Format("20060102")
	}

	return cfg
}

func (cmd *cmdEnv) checkNew(class *gradebook.Class, cfg newCfg) {
	if cmd.noOp() {
		return
	}

	isValidName(cmd, cfg.gbName)
	isValidType(cmd, cfg.gbType, class)
	isValidDate(cmd, cfg.gbDate)
}

func (cmd *cmdEnv) newGradebook(class *gradebook.Class, cfg newCfg) {
	if cmd.noOp() {
		return
	}

	emails := class.EmailsSortedByStudentName()
	recs := make(gradebook.AssignmentRecords, 0, len(emails))
	for _, email := range emails {
		recs = append(recs, &gradebook.AssignmentRecord{Email: email, Grade: nil})
	}

	newGb := &gradebook.Gradebook{
		AssignmentCategory: class.CategoriesByAssignmentType[cfg.gbType],
		AssignmentDate:     cfg.gbDate,
		AssignmentName:     cfg.gbName,
		AssignmentType:     cfg.gbType,
		AssignmentRecords:  recs,
	}

	gbData, err := json.MarshalIndent(newGb, "", "    ")
	if err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: problem marshaling gradebook: %s\n", cmd.name, err)

		return
	}

	fileName := fmt.Sprintf("%s-%s-%s.gradebook", cfg.gbType, cfg.gbName, cfg.gbDate)
	fileName = filepath.Join(cmd.directory, fileName)

	err = writeFile(fileName, gbData)
	if err != nil {
		cmd.exitValue = exitFailure
		if errors.Is(err, os.ErrExist) {
			fmt.Fprintf(cmd.stderr, "%s: %q already exists\n", cmd.name, fileName)
		} else {
			fmt.Fprintf(cmd.stderr, "%s: problem writing %q: %s\n", cmd.name, fileName, err)
		}
	}
}

func isValidName(cmd *cmdEnv, gbName string) {
	if gbName == "" || invalidGbNameRegex.MatchString(gbName) {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -name: %q\n", cmd.name, gbName)
	}
}

func isValidType(cmd *cmdEnv, gbType string, class *gradebook.Class) {
	if cmd.minNoOp() {
		return
	}

	gbTypes := slices.Collect(maps.Keys(class.CategoriesByAssignmentType))
	if !slices.Contains(gbTypes, gbType) {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -type: %q\n", cmd.name, gbType)
	}
}

func isValidDate(cmd *cmdEnv, gbDate string) {
	if cmd.minNoOp() {
		return
	}

	if _, err := time.Parse("20060102", gbDate); err != nil {
		cmd.exitValue = exitFailure
		fmt.Fprintf(cmd.stderr, "%s: invalid argument for -date: %q\n", cmd.name, gbDate)
	}
}

func writeFile(fileName string, data []byte) (err error) {
	fh, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("open file %q: %w", fileName, err)
	}
	defer func() {
		closeErr := fh.Close()
		if closeErr == nil {
			return
		}

		err = errors.Join(err, fmt.Errorf("close file %q: %w", fileName, closeErr))
	}()

	if _, err = fh.Write(data); err != nil {
		return fmt.Errorf("write file %q: %w", fileName, err)
	}

	syncErr := fh.Sync()
	if syncErr != nil {
		return fmt.Errorf("sync file %q: %w", fileName, syncErr)
	}

	return nil
}
