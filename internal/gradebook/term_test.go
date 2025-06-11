package gradebook_test

import (
	"testing"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

func TestTermIncludesValid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		term    *gradebook.Term
		dateStr string
		want    bool
	}{
		"term should include start date": {
			dateStr: "20230907",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should include end date": {
			dateStr: "20231011",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should include a date between start and end": {
			dateStr: "20230923",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: true,
		},
		"term should exclude a date before start": {
			dateStr: "20220907",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: false,
		},
		"term should exclude a date after end": {
			dateStr: "20231207",
			term: &gradebook.Term{
				Start: "20230907",
				End:   "20231011",
			},
			want: false,
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tc.term.Includes(tc.dateStr)
			if got != tc.want {
				t.Errorf("%#v.Includes(%s) returns %t; want %t", tc.term, tc.dateStr, got, tc.want)
			}
		})
	}
}

func TestLoadGradesWithTermFilter(t *testing.T) {
	t.Parallel()

	t.Run("loads grades within term", func(t *testing.T) {
		t.Parallel()

		class, err := gradebook.UnmarshalCalcClass("testdata/class.json")
		if err != nil {
			t.Fatalf("failed to unmarshal class: %v", err)
		}

		term := &gradebook.Term{
			Start: "20240301",
			End:   "20240331",
		}

		err = class.LoadGrades("testdata/term", term)
		if err != nil {
			t.Fatalf("LoadGrades failed: %v", err)
		}

		student := class.StudentsByEmail["gstriker@school.edu"]
		if len(student.GradesByCategory["minor"]) == 0 {
			t.Error("expected grades to be loaded for student within term")
		}
	})

	t.Run("excludes grades outside term", func(t *testing.T) {
		t.Parallel()

		class, err := gradebook.UnmarshalCalcClass("testdata/class.json")
		if err != nil {
			t.Fatalf("failed to unmarshal class: %v", err)
		}

		term := &gradebook.Term{
			Start: "20250101",
			End:   "20250131",
		}

		err = class.LoadGrades("testdata/term", term)
		if err != nil {
			t.Fatalf("LoadGrades failed: %v", err)
		}

		student := class.StudentsByEmail["gstriker@school.edu"]
		if len(student.GradesByCategory["minor"]) > 0 {
			t.Error("expected no grades to be loaded for term with no matching files")
		}
	})

	t.Run("nil term loads all grades", func(t *testing.T) {
		t.Parallel()

		class, err := gradebook.UnmarshalCalcClass("testdata/class.json")
		if err != nil {
			t.Fatalf("failed to unmarshal class: %v", err)
		}

		err = class.LoadGrades("testdata/term", nil)
		if err != nil {
			t.Fatalf("LoadGrades failed: %v", err)
		}

		student := class.StudentsByEmail["gstriker@school.edu"]
		if len(student.GradesByCategory["minor"]) == 0 {
			t.Error("expected grades to be loaded when term is nil")
		}
	})
}

func TestLoadGradesTermBoundaries(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		termStart    string
		termEnd      string
		expectedFile string
		shouldLoad   bool
	}{
		"march file on start date should be included": {
			termStart:    "20240315",
			termEnd:      "20240331",
			expectedFile: "march",
			shouldLoad:   true,
		},
		"march file on end date should be included": {
			termStart:    "20240301",
			termEnd:      "20240315",
			expectedFile: "march",
			shouldLoad:   true,
		},
		"march file before start should be excluded": {
			termStart:    "20240320",
			termEnd:      "20240331",
			expectedFile: "march",
			shouldLoad:   false,
		},
		"march file after end should be excluded": {
			termStart:    "20240301",
			termEnd:      "20240314",
			expectedFile: "march",
			shouldLoad:   false,
		},
		"april file on start date should be included": {
			termStart:    "20240415",
			termEnd:      "20240430",
			expectedFile: "april",
			shouldLoad:   true,
		},
		"april file on end date should be included": {
			termStart:    "20240401",
			termEnd:      "20240415",
			expectedFile: "april",
			shouldLoad:   true,
		},
		"april file before start should be excluded": {
			termStart:    "20240420",
			termEnd:      "20240430",
			expectedFile: "april",
			shouldLoad:   false,
		},
		"april file after end should be excluded": {
			termStart:    "20240401",
			termEnd:      "20240414",
			expectedFile: "april",
			shouldLoad:   false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class, err := gradebook.UnmarshalCalcClass("testdata/class.json")
			if err != nil {
				t.Fatalf("failed to unmarshal class: %v", err)
			}

			term := &gradebook.Term{
				Start: tt.termStart,
				End:   tt.termEnd,
			}

			err = class.LoadGrades("testdata/term", term)
			if err != nil {
				t.Fatalf("LoadGrades failed: %v", err)
			}

			student := class.StudentsByEmail["gstriker@school.edu"]
			var hasGrades bool

			if tt.expectedFile == "march" {
				hasGrades = len(student.GradesByCategory["minor"]) > 0
			} else if tt.expectedFile == "april" {
				hasGrades = len(student.GradesByCategory["major"]) > 0
			}

			if tt.shouldLoad && !hasGrades {
				t.Errorf("expected grades to be loaded for %s file in term %s to %s",
					tt.expectedFile, tt.termStart, tt.termEnd)
			}
			if !tt.shouldLoad && hasGrades {
				t.Errorf("expected no grades to be loaded for %s file outside term %s to %s",
					tt.expectedFile, tt.termStart, tt.termEnd)
			}
		})
	}
}
