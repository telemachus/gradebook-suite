package cli

import (
	"bytes"
	"testing"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

func TestIsValidName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name       string
		wantExit   int
		wantStderr bool
	}{
		"valid name with letters": {
			name:     "assignment",
			wantExit: exitSuccess,
		},
		"valid name with numbers": {
			name:     "quiz123",
			wantExit: exitSuccess,
		},
		"valid name with allowed punctuation": {
			name:     "test-quiz_2.final",
			wantExit: exitSuccess,
		},
		"valid name with all allowed chars": {
			name:     "Test_Quiz-2.Final123",
			wantExit: exitSuccess,
		},
		"empty string is invalid": {
			name:       "",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"space character is invalid": {
			name:       "test quiz",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"special characters are invalid": {
			name:       "test@quiz",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"unicode characters are invalid": {
			name:       "tëst",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"slash character is invalid": {
			name:       "test/quiz",
			wantExit:   exitFailure,
			wantStderr: true,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			cmd := &cmdEnv{
				name:      "testcmd",
				stderr:    &stderr,
				exitValue: exitSuccess,
			}

			isValidName(cmd, tt.name)

			if cmd.exitValue != tt.wantExit {
				t.Errorf("exitValue = %d; want %d", cmd.exitValue, tt.wantExit)
			}

			gotStderr := stderr.String()
			if !tt.wantStderr && gotStderr != "" {
				t.Errorf("expected empty stderr, got %q", gotStderr)
			}
			if tt.wantStderr && gotStderr == "" {
				t.Errorf("expected non-empty stderr, got empty")
			}
		})
	}
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

	tests := map[string]struct {
		gbType     string
		wantExit   int
		wantStderr bool
	}{
		"valid type": {
			gbType:   "test",
			wantExit: exitSuccess,
		},
		"another valid type": {
			gbType:   "quiz",
			wantExit: exitSuccess,
		},
		"third valid type": {
			gbType:   "cp",
			wantExit: exitSuccess,
		},
		"invalid type": {
			gbType:     "major",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"empty type": {
			gbType:     "",
			wantExit:   exitFailure,
			wantStderr: true,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			cmd := &cmdEnv{
				name:      "testcmd",
				stderr:    &stderr,
				exitValue: exitSuccess,
			}

			isValidType(cmd, tt.gbType, class)

			if cmd.exitValue != tt.wantExit {
				t.Errorf("exitValue = %d; want %d", cmd.exitValue, tt.wantExit)
			}

			gotStderr := stderr.String()
			if tt.wantStderr && gotStderr == "" {
				t.Error("expected non-empty stderr, got empty")
			}
			if !tt.wantStderr && gotStderr != "" {
				t.Errorf("expected empty stderr, got %q", gotStderr)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		gbDate     string
		wantExit   int
		wantStderr bool
	}{
		"valid date": {
			gbDate:   "20240319",
			wantExit: exitSuccess,
		},
		"another valid date": {
			gbDate:   "19991231",
			wantExit: exitSuccess,
		},
		"invalid date format - too short": {
			gbDate:     "2024319",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"invalid date format - too long": {
			gbDate:     "202403199",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"invalid date format - with dashes": {
			gbDate:     "2024-03-19",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"invalid date format - with slashes": {
			gbDate:     "03/19/2024",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"invalid date format - letters": {
			gbDate:     "202403ab",
			wantExit:   exitFailure,
			wantStderr: true,
		},
		"empty date": {
			gbDate:     "",
			wantExit:   exitFailure,
			wantStderr: true,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			cmd := &cmdEnv{
				name:      "testcmd",
				stderr:    &stderr,
				exitValue: exitSuccess,
			}

			isValidDate(cmd, tt.gbDate)

			if cmd.exitValue != tt.wantExit {
				t.Errorf("exitValue = %d; want %d", cmd.exitValue, tt.wantExit)
			}

			gotStderr := stderr.String()
			if tt.wantStderr && gotStderr == "" {
				t.Error("expected non-empty stderr, got empty")
			}
			if !tt.wantStderr && gotStderr != "" {
				t.Errorf("expected empty stderr, got %q", gotStderr)
			}
		})
	}
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

	tests := map[string]struct {
		name          string
		exitValue     int
		helpWanted    bool
		versionWanted bool
		cfg           newCfg
		wantExit      int
		wantStderr    bool
	}{
		"when exitValue is exitFailure, should not validate": {
			name:      "noOp due to exitFailure",
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
			name:       "noOp due to helpWanted",
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
			name:          "noOp due to versionWanted",
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
			name:      "validation should run",
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
			name:      "validation should run",
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
			name:      "validation should run",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			cmd := &cmdEnv{
				name:          "testcmd",
				stderr:        &stderr,
				exitValue:     tt.exitValue,
				helpWanted:    tt.helpWanted,
				versionWanted: tt.versionWanted,
			}

			cmd.checkNew(tt.cfg, class)

			if cmd.exitValue != tt.wantExit {
				t.Errorf("exitValue = %d; want %d", cmd.exitValue, tt.wantExit)
			}

			gotStderr := stderr.String()
			if !tt.wantStderr && gotStderr != "" {
				t.Errorf("expected empty stderr, got %q", gotStderr)
			}
			if tt.wantStderr && gotStderr == "" {
				t.Errorf("expected non-empty stderr, got empty")
			}
		})
	}
}
