package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/telemachus/gradebook"
)

var (
	stdioMu                       sync.Mutex
	errNoOutput                   = errors.New("no output")
	errCategoryLineBeforeStudent  = errors.New("category line before student line")
	errInvalidCategoryLineFormat  = errors.New("invalid category line format")
	errInvalidCategoryValueFormat = errors.New("invalid category value format")
)

func TestPublicGradebookNames(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookNames, []string{"-dir", dir})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}

	want := "Bob Young\nAlice Zephyr\n"
	if stdout != want {
		t.Fatalf("stdout mismatch:\nwant:\n%q\ngot:\n%q", want, stdout)
	}
}

func TestPublicGradebookNamesLastFirst(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookNames, []string{"-dir", dir, "-last-first"})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}

	want := "Young, Bob\nZephyr, Alice\n"
	if stdout != want {
		t.Fatalf("stdout mismatch:\nwant:\n%q\ngot:\n%q", want, stdout)
	}
}

func TestPublicGradebookEmails(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookEmails, []string{"-dir", dir})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}

	want := "bob@example.com\nalice@example.com\n"
	if stdout != want {
		t.Fatalf("stdout mismatch:\nwant:\n%q\ngot:\n%q", want, stdout)
	}
}

func TestPublicGradebookCalc(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookCalc, []string{"-dir", dir, "-term", "q1"})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}

	want := "" +
		"Bob Young\n" +
		"\tOverall average: 90\n" +
		"\tMajor: No results\n" +
		"\tMinor: 90\n" +
		"\tParticipation: No results\n" +
		"Alice Zephyr\n" +
		"\tOverall average: No results\n" +
		"\tMajor: No results\n" +
		"\tMinor: No results\n" +
		"\tParticipation: No results\n"
	if stdout != want {
		t.Fatalf("stdout mismatch:\nwant:\n%q\ngot:\n%q", want, stdout)
	}
}

func TestPublicGradebookUnscored(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookUnscored, []string{"-dir", dir, "-term", "q1"})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}

	got, err := parseUnscoredOutput(stdout)
	if err != nil {
		t.Fatalf("parseUnscoredOutput() returned error: %v", err)
	}

	want := map[string]map[string]int{
		"Bob Young": {
			"Major":         0,
			"Minor":         0,
			"Participation": 0,
		},
		"Alice Zephyr": {
			"Major":         0,
			"Minor":         1,
			"Participation": 0,
		},
	}

	assertUnscoredCounts(t, got, want)
}

func TestPublicGradebookNewCreatesFile(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)

	exitCode, stdout, stderr := runPublicCommand(t, GradebookNew, []string{
		"-dir", dir,
		"-name", "unit",
		"-type", "quiz",
		"-date", "20240501",
	})
	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q; want empty", stdout)
	}

	path := filepath.Join(dir, "quiz-unit-20240501.gradebook")
	gb, err := gradebook.UnmarshalGradebook(path)
	if err != nil {
		t.Fatalf("failed to unmarshal created gradebook: %v", err)
	}
	if gb.AssignmentType != "quiz" {
		t.Fatalf("AssignmentType = %q; want %q", gb.AssignmentType, "quiz")
	}
}

func TestPublicGradebookCalcHelpShortCircuit(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runPublicCommand(t, GradebookCalc, []string{"-help", "-dir", "/path/that/does/not/exist"})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if !strings.Contains(stdout, "usage: gradebook-calc") {
		t.Fatalf("stdout = %q; expected help usage output", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}
}

func TestPublicGradebookNamesVersionShortCircuit(t *testing.T) {
	t.Parallel()

	exitCode, stdout, stderr := runPublicCommand(t, GradebookNames, []string{"-version", "-dir", "/path/that/does/not/exist"})

	if exitCode != exitSuccess {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitSuccess)
	}
	if !strings.Contains(stdout, "gradebook-names: "+suiteVersion) {
		t.Fatalf("stdout = %q; expected version output", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q; want empty", stderr)
	}
}

func TestPublicGradebookCalcInvalidTerm(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookCalc, []string{"-dir", dir, "-term", "not-a-term"})

	if exitCode != exitFailure {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitFailure)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q; want empty", stdout)
	}
	wantStderr := "gradebook-calc: \"not-a-term\" is not a valid term\n"
	if stderr != wantStderr {
		t.Fatalf("stderr mismatch:\nwant:\n%q\ngot:\n%q", wantStderr, stderr)
	}
}

func TestPublicGradebookNewInvalidNameFailFast(t *testing.T) {
	t.Parallel()

	dir := writeSuiteFixture(t)
	exitCode, stdout, stderr := runPublicCommand(t, GradebookNew, []string{
		"-dir", dir,
		"-name", "bad name",
		"-type", "quiz",
		"-date", "20240501",
	})

	if exitCode != exitFailure {
		t.Fatalf("exitCode = %d; want %d", exitCode, exitFailure)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q; want empty", stdout)
	}
	wantStderr := "gradebook-new: invalid argument for -name: \"bad name\"\n"
	if stderr != wantStderr {
		t.Fatalf("stderr mismatch:\nwant:\n%q\ngot:\n%q", wantStderr, stderr)
	}

	createdPath := filepath.Join(dir, "quiz-bad name-20240501.gradebook")
	if _, err := os.Stat(createdPath); !os.IsNotExist(err) {
		t.Fatalf("expected no file at %q; stat err = %v", createdPath, err)
	}
}

func parseUnscoredOutput(stdout string) (map[string]map[string]int, error) {
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) == 0 {
		return nil, errNoOutput
	}

	result := make(map[string]map[string]int)
	student := ""
	for _, line := range lines {
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "\t") {
			student = strings.TrimSuffix(line, ":")
			result[student] = make(map[string]int)

			continue
		}

		if student == "" {
			return nil, fmt.Errorf("%w: %q", errCategoryLineBeforeStudent, line)
		}

		trimmed := strings.TrimSpace(line)
		parts := strings.SplitN(trimmed, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%w: %q", errInvalidCategoryLineFormat, line)
		}

		label := parts[0]
		words := strings.Split(parts[1], " ")
		if len(words) < 3 {
			return nil, fmt.Errorf("%w: %q", errInvalidCategoryValueFormat, line)
		}

		count, err := strconv.Atoi(words[0])
		if err != nil {
			return nil, fmt.Errorf("invalid unscored count in line %q: %w", line, err)
		}

		result[student][label] = count
	}

	return result, nil
}

func assertUnscoredCounts(t *testing.T, got, want map[string]map[string]int) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d; want %d", len(got), len(want))
	}

	for student, wantByCategory := range want {
		gotByCategory, ok := got[student]
		if !ok {
			t.Fatalf("student %q missing from parsed output", student)
		}

		if len(gotByCategory) != len(wantByCategory) {
			t.Fatalf(
				"len(got[%q]) = %d; want %d",
				student,
				len(gotByCategory),
				len(wantByCategory),
			)
		}

		for category, wantCount := range wantByCategory {
			gotCount, ok := gotByCategory[category]
			if !ok {
				t.Fatalf("category %q missing for student %q", category, student)
			}

			if gotCount != wantCount {
				t.Fatalf(
					"got[%q][%q] = %d; want %d",
					student,
					category,
					gotCount,
					wantCount,
				)
			}
		}
	}
}

func runPublicCommand(t *testing.T, cmdFunc func([]string) int, args []string) (int, string, string) {
	t.Helper()

	stdioMu.Lock()
	defer stdioMu.Unlock()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed creating stdout pipe: %v", err)
	}
	defer func() {
		if closeErr := stdoutR.Close(); closeErr != nil {
			t.Fatalf("failed closing stdout reader: %v", closeErr)
		}
	}()

	stderrR, stderrW, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed creating stderr pipe: %v", err)
	}
	defer func() {
		if closeErr := stderrR.Close(); closeErr != nil {
			t.Fatalf("failed closing stderr reader: %v", closeErr)
		}
	}()

	os.Stdout = stdoutW
	os.Stderr = stderrW

	exitCode := cmdFunc(args)

	if err = stdoutW.Close(); err != nil {
		t.Fatalf("failed closing stdout writer: %v", err)
	}
	if err = stderrW.Close(); err != nil {
		t.Fatalf("failed closing stderr writer: %v", err)
	}

	stdoutData, err := io.ReadAll(stdoutR)
	if err != nil {
		t.Fatalf("failed reading stdout: %v", err)
	}
	stderrData, err := io.ReadAll(stderrR)
	if err != nil {
		t.Fatalf("failed reading stderr: %v", err)
	}

	return exitCode, string(stdoutData), string(stderrData)
}
