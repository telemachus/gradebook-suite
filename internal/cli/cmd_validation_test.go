package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validClassForCmdValidation = `{
    "name": "Validation Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice@example.com": {
            "first_name": "Alice",
            "last_name": "Zephyr"
        }
    }
}`

const invalidTermClassForCmdValidation = `{
    "name": "Validation Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240230",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice@example.com": {
            "first_name": "Alice",
            "last_name": "Zephyr"
        }
    }
}`

const nilStudentClassForCmdValidation = `{
    "name": "Validation Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice@example.com": null
    }
}`

const badEmailClassForCmdValidation = `{
    "name": "Validation Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice.example.com": {
            "first_name": "Alice",
            "last_name": "Zephyr"
        }
    }
}`

const emptyFirstNameClassForCmdValidation = `{
    "name": "Validation Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major"],
    "labels_by_assignment_category": {
        "major": "Major"
    },
    "weights_by_assignment_category": {
        "major": 100
    },
    "categories_by_assignment_type": {
        "test": "major"
    },
    "students_by_email": {
        "alice@example.com": {
            "first_name": "",
            "last_name": "Zephyr"
        }
    }
}`

func TestUnmarshalClassValidatesClass(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"invalid term date":        invalidTermClassForCmdValidation,
		"nil student pointer":      nilStudentClassForCmdValidation,
		"student email missing @":  badEmailClassForCmdValidation,
		"empty student first name": emptyFirstNameClassForCmdValidation,
	}

	for testName, classData := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			classFile := writeClassFixture(t, classData)
			cmd, stderr := newValidationCmd(classFile)
			class := cmd.unmarshalClass()

			if class != nil {
				t.Fatal("unmarshalClass() returned non-nil class for invalid class fixture")
			}
			if cmd.exitValue != exitFailure {
				t.Fatalf("exitValue = %d; want %d", cmd.exitValue, exitFailure)
			}
			if !strings.Contains(stderr.String(), "problem validating class:") {
				t.Fatalf("stderr = %q; want validation error message", stderr.String())
			}
		})
	}
}

func TestUnmarshalClassValidationPassesForValidClass(t *testing.T) {
	t.Parallel()

	classFile := writeClassFixture(t, validClassForCmdValidation)
	cmd, stderr := newValidationCmd(classFile)
	class := cmd.unmarshalClass()

	if class == nil {
		t.Fatal("unmarshalClass() returned nil for valid class fixture")
	}
	if cmd.exitValue != exitSuccess {
		t.Fatalf("exitValue = %d; want %d", cmd.exitValue, exitSuccess)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q; want empty", got)
	}
}

func writeClassFixture(t *testing.T, classData string) string {
	t.Helper()

	classFile := filepath.Join(t.TempDir(), "class.json")
	if err := os.WriteFile(classFile, []byte(classData), 0o644); err != nil {
		t.Fatalf("failed writing class fixture: %v", err)
	}

	return classFile
}

func newValidationCmd(classFile string) (*cmdEnv, *bytes.Buffer) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := cmdFromWithWriters("gradebook-names", namesUsage, &stdout, &stderr)
	cmd.classFile = classFile

	return cmd, &stderr
}
