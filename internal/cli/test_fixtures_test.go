package cli

import (
	"os"
	"path/filepath"
	"testing"
)

const classFixtureJSON = `{
    "name": "Characterization Test Class",
    "terms_by_id": {
        "q1": {
            "start": "20240101",
            "end": "20241231"
        }
    },
    "assignment_categories": ["major", "minor", "cp"],
    "labels_by_assignment_category": {
        "major": "Major",
        "minor": "Minor",
        "cp": "Participation"
    },
    "weights_by_assignment_category": {
        "major": 50,
        "minor": 30,
        "cp": 20
    },
    "categories_by_assignment_type": {
        "quiz": "minor",
        "test": "major",
        "cp": "cp"
    },
    "students_by_email": {
        "alice@example.com": {
            "first_name": "Alice",
            "last_name": "Zephyr"
        },
        "bob@example.com": {
            "first_name": "Bob",
            "last_name": "Young"
        }
    }
}`

const gradebookFixtureJSON = `{
    "assignment_category": "minor",
    "assignment_date": "20240319",
    "assignment_records": [
        {
            "email": "bob@example.com",
            "grade": 90
        },
        {
            "email": "alice@example.com",
            "grade": null
        }
    ],
    "assignment_name": "quiz-1",
    "assignment_type": "quiz"
}`

func writeSuiteFixture(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	mustWriteFixtureFile(t, filepath.Join(dir, suiteClassFile), classFixtureJSON)
	mustWriteFixtureFile(t, filepath.Join(dir, "quiz-quiz-1-20240319.gradebook"), gradebookFixtureJSON)

	return dir
}

func mustWriteFixtureFile(t *testing.T, path, data string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("failed writing fixture file %q: %v", path, err)
	}
}
