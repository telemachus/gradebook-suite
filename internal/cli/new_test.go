package cli

import (
	"bytes"
	"testing"

	"github.com/telemachus/gradebook"
)

type validationCase struct {
	input      string
	wantExit   int
	wantStderr bool
}

var isValidNameCases = map[string]validationCase{
	"valid name with letters": {
		input:    "assignment",
		wantExit: exitSuccess,
	},
	"valid name with numbers": {
		input:    "quiz123",
		wantExit: exitSuccess,
	},
	"valid name with allowed punctuation": {
		input:    "test-quiz_2.final",
		wantExit: exitSuccess,
	},
	"valid name with all allowed chars": {
		input:    "Test_Quiz-2.Final123",
		wantExit: exitSuccess,
	},
	"empty string is invalid": {
		input:      "",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"space character is invalid": {
		input:      "test quiz",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"special characters are invalid": {
		input:      "test@quiz",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"unicode characters are invalid": {
		input:      "tëst",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"slash character is invalid": {
		input:      "test/quiz",
		wantExit:   exitFailure,
		wantStderr: true,
	},
}

func TestIsValidName(t *testing.T) {
	t.Parallel()

	for testName, tt := range isValidNameCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd, stderr := newTestCmdEnv()
			isValidName(cmd, tt.input)
			assertCmdResult(t, cmd, stderr, tt.wantExit, tt.wantStderr)
		})
	}
}

var isValidTypeCases = map[string]validationCase{
	"valid type": {
		input:    "test",
		wantExit: exitSuccess,
	},
	"another valid type": {
		input:    "quiz",
		wantExit: exitSuccess,
	},
	"third valid type": {
		input:    "cp",
		wantExit: exitSuccess,
	},
	"invalid type": {
		input:      "major",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"empty type": {
		input:      "",
		wantExit:   exitFailure,
		wantStderr: true,
	},
}

func TestIsValidType(t *testing.T) {
	t.Parallel()

	class := &gradebook.Class{
		CategoriesByAssignmentType: gradebook.CategoriesByAssignmentType{
			"test": "major",
			"quiz": "minor",
			"cp":   "cp",
		},
	}

	for testName, tt := range isValidTypeCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd, stderr := newTestCmdEnv()
			isValidType(cmd, tt.input, class)
			assertCmdResult(t, cmd, stderr, tt.wantExit, tt.wantStderr)
		})
	}
}

var isValidDateCases = map[string]validationCase{
	"valid date": {
		input:    "20240319",
		wantExit: exitSuccess,
	},
	"another valid date": {
		input:    "19991231",
		wantExit: exitSuccess,
	},
	"invalid date format - too short": {
		input:      "2024319",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"invalid date format - too long": {
		input:      "202403199",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"invalid date format - with dashes": {
		input:      "2024-03-19",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"invalid date format - with slashes": {
		input:      "03/19/2024",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"invalid date format - letters": {
		input:      "202403ab",
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"empty date": {
		input:      "",
		wantExit:   exitFailure,
		wantStderr: true,
	},
}

func TestIsValidDate(t *testing.T) {
	t.Parallel()

	for testName, tt := range isValidDateCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd, stderr := newTestCmdEnv()
			isValidDate(cmd, tt.input)
			assertCmdResult(t, cmd, stderr, tt.wantExit, tt.wantStderr)
		})
	}
}

type checkNewCase struct {
	cfg           newCfg
	exitValue     int
	wantExit      int
	helpWanted    bool
	versionWanted bool
	wantStderr    bool
}

var checkNewCases = map[string]checkNewCase{
	"when exitValue is exitFailure, should not validate": {
		exitValue: exitFailure,
		cfg: newCfg{
			gbName: "invalid@name",
			gbType: "invalidtype",
			gbDate: "invaliddate",
		},
		wantExit:   exitFailure,
		wantStderr: false,
	},
	"when helpWanted is true, should not validate": {
		exitValue:  exitSuccess,
		helpWanted: true,
		cfg: newCfg{
			gbName: "invalid@name",
			gbType: "invalidtype",
			gbDate: "invaliddate",
		},
		wantExit:   exitSuccess,
		wantStderr: false,
	},
	"when versionWanted is true, should not validate": {
		exitValue:     exitSuccess,
		versionWanted: true,
		cfg: newCfg{
			gbName: "invalid@name",
			gbType: "invalidtype",
			gbDate: "invaliddate",
		},
		wantExit:   exitSuccess,
		wantStderr: false,
	},
	"when noOp is false, should validate name and fail": {
		exitValue: exitSuccess,
		cfg: newCfg{
			gbName: "invalid@name",
			gbType: "quiz",
			gbDate: "20240319",
		},
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"when noOp is false, should validate type and fail": {
		exitValue: exitSuccess,
		cfg: newCfg{
			gbName: "name",
			gbType: "major",
			gbDate: "20240319",
		},
		wantExit:   exitFailure,
		wantStderr: true,
	},
	"when noOp is false, should validate date and fail": {
		exitValue: exitSuccess,
		cfg: newCfg{
			gbName: "name",
			gbType: "major",
			gbDate: "2024031",
		},
		wantExit:   exitFailure,
		wantStderr: true,
	},
}

func TestCheckNew(t *testing.T) {
	t.Parallel()

	class := &gradebook.Class{
		CategoriesByAssignmentType: gradebook.CategoriesByAssignmentType{
			"test": "major",
			"quiz": "minor",
			"cp":   "cp",
		},
	}

	for testName, tt := range checkNewCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			cmd, stderr := newTestCmdEnv()
			cmd.exitValue = tt.exitValue
			cmd.helpWanted = tt.helpWanted
			cmd.versionWanted = tt.versionWanted

			cmd.checkNew(class, tt.cfg)
			assertCmdResult(t, cmd, stderr, tt.wantExit, tt.wantStderr)
		})
	}
}

func newTestCmdEnv() (*cmdEnv, *bytes.Buffer) {
	var stderr bytes.Buffer
	cmd := &cmdEnv{
		name:      "testcmd",
		stderr:    &stderr,
		exitValue: exitSuccess,
	}

	return cmd, &stderr
}

func assertCmdResult(t *testing.T, cmd *cmdEnv, stderr *bytes.Buffer, wantExit int, wantStderr bool) {
	t.Helper()

	if cmd.exitValue != wantExit {
		t.Errorf("exitValue = %d; want %d", cmd.exitValue, wantExit)
	}

	gotStderr := stderr.String()
	if !wantStderr && gotStderr != "" {
		t.Errorf("expected empty stderr, got %q", gotStderr)
	}
	if wantStderr && gotStderr == "" {
		t.Errorf("expected non-empty stderr, got empty")
	}
}
