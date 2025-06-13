package gradebook_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

func TestStudentsSortedByName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		students map[string]*gradebook.Student
		want     []*gradebook.Student
	}{
		"empty class returns empty slice": {
			students: map[string]*gradebook.Student{},
			want:     []*gradebook.Student{},
		},
		"single student returns single student": {
			students: map[string]*gradebook.Student{
				"alice@example.com": {
					FirstName: "Alice",
					LastName:  "Anderson",
				},
			},
			want: []*gradebook.Student{
				{
					FirstName: "Alice",
					LastName:  "Anderson",
				},
			},
		},
		"multiple students sorted by last name": {
			students: map[string]*gradebook.Student{
				"charlie@example.com": {
					FirstName: "Charlie",
					LastName:  "Chen",
				},
				"alice@example.com": {
					FirstName: "Alice",
					LastName:  "Anderson",
				},
				"bob@example.com": {
					FirstName: "Bob",
					LastName:  "Baker",
				},
			},
			want: []*gradebook.Student{
				{
					FirstName: "Alice",
					LastName:  "Anderson",
				},
				{
					FirstName: "Bob",
					LastName:  "Baker",
				},
				{
					FirstName: "Charlie",
					LastName:  "Chen",
				},
			},
		},
		"same last name sorted by first name": {
			students: map[string]*gradebook.Student{
				"bob.smith@example.com": {
					FirstName: "Bob",
					LastName:  "Smith",
				},
				"alice.smith@example.com": {
					FirstName: "Alice",
					LastName:  "Smith",
				},
				"charlie.smith@example.com": {
					FirstName: "Charlie",
					LastName:  "Smith",
				},
			},
			want: []*gradebook.Student{
				{
					FirstName: "Alice",
					LastName:  "Smith",
				},
				{
					FirstName: "Bob",
					LastName:  "Smith",
				},
				{
					FirstName: "Charlie",
					LastName:  "Smith",
				},
			},
		},
		"mixed sorting - last name priority, then first name": {
			students: map[string]*gradebook.Student{
				"bob.young@example.com": {
					FirstName: "Bob",
					LastName:  "Young",
				},
				"alice.young@example.com": {
					FirstName: "Alice",
					LastName:  "Young",
				},
				"charlie.smith@example.com": {
					FirstName: "Charlie",
					LastName:  "Smith",
				},
				"david.anderson@example.com": {
					FirstName: "David",
					LastName:  "Anderson",
				},
			},
			want: []*gradebook.Student{
				{
					FirstName: "David",
					LastName:  "Anderson",
				},
				{
					FirstName: "Charlie",
					LastName:  "Smith",
				},
				{
					FirstName: "Alice",
					LastName:  "Young",
				},
				{
					FirstName: "Bob",
					LastName:  "Young",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &gradebook.Class{
				StudentsByEmail: tt.students,
			}

			got := class.StudentsSortedByName()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("StudentsSortedByName() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEmailsSortedByStudentName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		students map[string]*gradebook.Student
		want     []string
	}{
		"empty class returns empty slice": {
			students: map[string]*gradebook.Student{},
			want:     []string{},
		},
		"single student returns single email": {
			students: map[string]*gradebook.Student{
				"alice@example.com": {
					FirstName: "Alice",
					LastName:  "Anderson",
				},
			},
			want: []string{"alice@example.com"},
		},
		"multiple students sorted by last name": {
			students: map[string]*gradebook.Student{
				"charlie@example.com": {
					FirstName: "Charlie",
					LastName:  "Chen",
				},
				"alice@example.com": {
					FirstName: "Alice",
					LastName:  "Anderson",
				},
				"bob@example.com": {
					FirstName: "Bob",
					LastName:  "Baker",
				},
			},
			want: []string{
				"alice@example.com",
				"bob@example.com",
				"charlie@example.com",
			},
		},
		"same last name sorted by first name": {
			students: map[string]*gradebook.Student{
				"bob.smith@example.com": {
					FirstName: "Bob",
					LastName:  "Smith",
				},
				"alice.smith@example.com": {
					FirstName: "Alice",
					LastName:  "Smith",
				},
				"charlie.smith@example.com": {
					FirstName: "Charlie",
					LastName:  "Smith",
				},
			},
			want: []string{
				"alice.smith@example.com",
				"bob.smith@example.com",
				"charlie.smith@example.com",
			},
		},
		"mixed sorting - last name priority, then first name": {
			students: map[string]*gradebook.Student{
				"bob.young@example.com": {
					FirstName: "Bob",
					LastName:  "Young",
				},
				"alice.young@example.com": {
					FirstName: "Alice",
					LastName:  "Young",
				},
				"charlie.smith@example.com": {
					FirstName: "Charlie",
					LastName:  "Smith",
				},
				"david.anderson@example.com": {
					FirstName: "David",
					LastName:  "Anderson",
				},
			},
			want: []string{
				"david.anderson@example.com",
				"charlie.smith@example.com",
				"alice.young@example.com",
				"bob.young@example.com",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &gradebook.Class{
				StudentsByEmail: tt.students,
			}

			got := class.EmailsSortedByStudentName()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("EmailsSortedByStudentName() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAssignmentCategoriesSortedByLabel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		categories gradebook.AssignmentCategories
		labels     gradebook.LabelsByAssignmentCategory
		want       []string
	}{
		"empty categories returns empty slice": {
			categories: gradebook.AssignmentCategories{},
			labels:     gradebook.LabelsByAssignmentCategory{},
			want:       []string{},
		},
		"single category returns single category": {
			categories: gradebook.AssignmentCategories{"major"},
			labels: gradebook.LabelsByAssignmentCategory{
				"major": "Major Assessments",
			},
			want: []string{"major"},
		},
		"multiple categories sorted by label": {
			categories: gradebook.AssignmentCategories{"major", "cp", "minor"},
			labels: gradebook.LabelsByAssignmentCategory{
				"major": "Major Assessments",
				"minor": "Daily Work",
				"cp":    "Class Participation",
			},
			want: []string{"cp", "minor", "major"},
		},
		"categories with same label prefix sorted alphabetically": {
			categories: gradebook.AssignmentCategories{"quiz", "exam", "project"},
			labels: gradebook.LabelsByAssignmentCategory{
				"quiz":    "Assessment: Quiz",
				"exam":    "Assessment: Exam",
				"project": "Assessment: Project",
			},
			want: []string{"exam", "project", "quiz"},
		},
		"mixed label sorting": {
			categories: gradebook.AssignmentCategories{"final", "hw", "participation", "midterm"},
			labels: gradebook.LabelsByAssignmentCategory{
				"final":         "Final Exam",
				"hw":            "Homework",
				"participation": "Class Participation",
				"midterm":       "Midterm Exam",
			},
			want: []string{"participation", "final", "hw", "midterm"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			class := &gradebook.Class{
				AssignmentCategories:       tt.categories,
				LabelsByAssignmentCategory: tt.labels,
			}

			got := class.AssignmentCategoriesSortedByLabel()

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AssignmentCategoriesSortedByLabel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
