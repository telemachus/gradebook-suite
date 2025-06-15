package gradebook_test

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

var (
	classJSON            = "testdata/class.json"
	classUnequalJSON     = "testdata/wrong.json"
	classInvalidJSON     = "testdata/invalid.json"
	gradebookJSON        = "testdata/quiz-golden-20240319.gradebook"
	gradebookUnequalJSON = "testdata/quiz-wrong-20240319.gradebook"
	gradebookInvalidJSON = "testdata/quiz-invalid-20240319.gradebook"
)

func TestUnmarshalClass(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := gradebook.UnmarshalClass(classJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalClass(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestUnmarshalClassUnequal(t *testing.T) {
	t.Parallel()

	want := fakeClass()
	got, err := gradebook.UnmarshalClass(classUnequalJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("gradebook.UnmarshalClass(classUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestUnmarshalClassInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.UnmarshalClass(classInvalidJSON)
	if err == nil {
		t.Fatal("want error; got nil")
	}
}

func TestUnmarshalGradebook(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := gradebook.UnmarshalGradebook(gradebookJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalGradebook(gradebookJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalGradebook(gradebookJSON) mismatch(-want +got):\n%s", diff)
	}
}

func TestUnmarshalGradebookUnequal(t *testing.T) {
	t.Parallel()

	want := fakeGradebook()
	got, err := gradebook.UnmarshalGradebook(gradebookUnequalJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalGradebook(gradebookUnequalJSON) = %v; want nil error", err)
	}

	if cmp.Equal(want, got) {
		t.Error("gradebook.UnmarshalGradebook(gradebookUnequalJSON) should differ from the mock class, but it does not")
	}
}

func TestUnmarshalGradebookInvalid(t *testing.T) {
	t.Parallel()

	_, err := gradebook.UnmarshalGradebook(gradebookInvalidJSON)
	if err == nil {
		t.Fatalf("gradebook.UnmarshalGradebook(%q) should error; got nil", gradebookInvalidJSON)
	}
}

func TestLoadGradesValid(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "validgradebooks")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err != nil {
		t.Fatalf("want no error from LoadGrades(%q, nil); got %v", testDir, err)
	}
}

func TestLoadGradesInvalid(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "invalidgradebook")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func TestLoadGradesUnknownType(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "unknowntypegradebook")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func TestLoadGradesNonexistentDirectory(t *testing.T) {
	t.Parallel()

	pseudoClass := fakeCalcClass()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("problem with os.Getwd(): cannot continue")
	}
	testDir := filepath.Join(wd, "testdata", "nonexistent")

	err = pseudoClass.LoadGrades(testDir, nil)
	if err == nil {
		t.Fatalf("want error from LoadGrades(%q, nil); got nil", testDir)
	}
}

func fakeCalcClass() *gradebook.Class {
	return &gradebook.Class{
		Name: "Lucretius",
		TermsByID: map[string]*gradebook.Term{
			"q1": {
				Start: "20200910",
				End:   "20201103",
			},
			"q2": {
				Start: "20201108",
				End:   "20210114",
			},
			"q3": {
				Start: "20210124",
				End:   "20210311",
			},
			"q4": {
				Start: "20210328",
				End:   "20210609",
			},
			"s1": {
				Start: "20200910",
				End:   "20210114",
			},
			"s2": {
				Start: "20210124",
				End:   "20210609",
			},
		},
		AssignmentCategories: gradebook.AssignmentCategories{"major", "minor", "cp"},
		LabelsByAssignmentCategory: gradebook.LabelsByAssignmentCategory{
			"major": "Major assessments",
			"minor": "Daily work and quizzes",
			"cp":    "Class participation",
		},
		WeightsByAssignmentCategory: gradebook.WeightsByAssignmentCategory{
			"major": 30,
			"minor": 40,
			"cp":    30,
		},
		CategoriesByAssignmentType: gradebook.CategoriesByAssignmentType{
			"test":    "major",
			"project": "major",
			"essay":   "major",
			"quiz":    "minor",
			"hw":      "minor",
			"cp":      "cp",
		},
		StudentsByEmail: gradebook.StudentsByEmail{
			"gstriker@school.edu": &gradebook.Student{
				FirstName:        "Gisela",
				LastName:         "Striker",
				GradesByCategory: fakeGradesMap(),
			},
			"mfrede@school.edu": &gradebook.Student{
				FirstName:        "Michael",
				LastName:         "Frede",
				GradesByCategory: fakeGradesMap(),
			},
			"jannas@school.edu": &gradebook.Student{
				FirstName:        "Julia",
				LastName:         "Annas",
				GradesByCategory: fakeGradesMap(),
			},
			"agomezlobo@school.edu": &gradebook.Student{
				FirstName:        "Alfonso",
				LastName:         "Gómez-Lobo",
				GradesByCategory: fakeGradesMap(),
			},
			"gfine@school.edu": &gradebook.Student{
				FirstName:        "Gail",
				LastName:         "Fine",
				GradesByCategory: fakeGradesMap(),
			},
		},
	}
}

func fakeGradebook() *gradebook.Gradebook {
	return &gradebook.Gradebook{
		AssignmentCategory: "minor",
		AssignmentDate:     "20240319",
		AssignmentName:     "golden",
		AssignmentType:     "quiz",
		AssignmentGrades: gradebook.Grades{
			&gradebook.Grade{
				Email: "gstriker@school.edu",
				Score: floatPtr(94.2),
			},
			&gradebook.Grade{
				Email: "mfrede@school.edu",
				Score: floatPtr(94.0),
			},
			&gradebook.Grade{
				Email: "jannas@school.edu",
				Score: floatPtr(104),
			},
			&gradebook.Grade{
				Email: "agomezlobo@school.edu",
				Score: floatPtr(81),
			},
			&gradebook.Grade{
				Email: "gfine@school.edu",
			},
		},
	}
}

func fakeClass() *gradebook.Class {
	return &gradebook.Class{
		Name: "Lucretius",
		TermsByID: map[string]*gradebook.Term{
			"q1": {
				Start: "20200910",
				End:   "20201103",
			},
			"q2": {
				Start: "20201108",
				End:   "20210114",
			},
			"q3": {
				Start: "20210124",
				End:   "20210311",
			},
			"q4": {
				Start: "20210328",
				End:   "20210609",
			},
			"s1": {
				Start: "20200910",
				End:   "20210114",
			},
			"s2": {
				Start: "20210124",
				End:   "20210609",
			},
		},
		AssignmentCategories: gradebook.AssignmentCategories{"major", "minor", "cp"},
		LabelsByAssignmentCategory: gradebook.LabelsByAssignmentCategory{
			"major": "Major assessments",
			"minor": "Daily work and quizzes",
			"cp":    "Class participation",
		},
		WeightsByAssignmentCategory: gradebook.WeightsByAssignmentCategory{
			"major": 30,
			"minor": 40,
			"cp":    30,
		},
		CategoriesByAssignmentType: gradebook.CategoriesByAssignmentType{
			"test":    "major",
			"project": "major",
			"essay":   "major",
			"quiz":    "minor",
			"hw":      "minor",
			"cp":      "cp",
		},
		StudentsByEmail: gradebook.StudentsByEmail{
			"gstriker@school.edu": &gradebook.Student{
				FirstName: "Gisela",
				LastName:  "Striker",
			},
			"mfrede@school.edu": &gradebook.Student{
				FirstName: "Michael",
				LastName:  "Frede",
			},
			"jannas@school.edu": &gradebook.Student{
				FirstName: "Julia",
				LastName:  "Annas",
			},
			"agomezlobo@school.edu": &gradebook.Student{
				FirstName: "Alfonso",
				LastName:  "Gómez-Lobo",
			},
			"gfine@school.edu": &gradebook.Student{
				FirstName: "Gail",
				LastName:  "Fine",
			},
		},
	}
}

func TestStudentAverage(t *testing.T) {
	t.Parallel()

	student, err := gradebook.NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.GradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
	}

	result := student.Average("major")
	if result.Valid {
		t.Error("Average(major) with no grades should return Valid=false")
	}

	grades := []float64{85, 90, 95}
	student.GradesByCategory["major"] = append(student.GradesByCategory["major"], grades...)

	result = student.Average("major")
	if !result.Valid {
		t.Error("Average(major) with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("Average(major) = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentAverageInvalidCategory(t *testing.T) {
	t.Parallel()

	student, err := gradebook.NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.GradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
	}
}

func TestStudentTotalAverage(t *testing.T) {
	t.Parallel()

	student, err := gradebook.NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.GradesByCategory = map[string][]float64{
		"major": make([]float64, 0),
		"minor": make([]float64, 0),
		"cp":    make([]float64, 0),
	}

	weights := gradebook.WeightsByAssignmentCategory{
		"major": 50,
		"minor": 30,
		"cp":    20,
	}

	result := student.TotalAverage(weights)
	if result.Valid {
		t.Error("TotalAverage() with no grades should return Valid=false")
	}

	student.GradesByCategory["major"] = append(student.GradesByCategory["major"], 90)
	student.GradesByCategory["minor"] = append(student.GradesByCategory["minor"], 90)
	student.GradesByCategory["cp"] = append(student.GradesByCategory["cp"], 90)

	result = student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("TotalAverage() with equal grades = %f; want %f", result.Value, expectedAvg)
	}

	student.GradesByCategory = map[string][]float64{
		"major": {94},
		"minor": {82},
		"cp":    {75},
	}

	result = student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with grades should return Valid=true")
	}

	expectedAvg = 86.6
	if !floatEqual(result.Value, expectedAvg, 0.1) {
		t.Errorf("TotalAverage() with different grades = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentTotalAveragePartialGrades(t *testing.T) {
	t.Parallel()

	student, err := gradebook.NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.GradesByCategory = map[string][]float64{
		"major": {90},
		"minor": {80},
		"cp":    make([]float64, 0),
	}

	weights := map[string]int{
		"major": 50,
		"minor": 30,
		"cp":    20,
	}

	result := student.TotalAverage(weights)
	if !result.Valid {
		t.Error("TotalAverage() with partial grades should return Valid=true")
	}

	expectedAvg := 86.25
	if !floatEqual(result.Value, expectedAvg, 0.01) {
		t.Errorf("TotalAverage() with partial grades = %f; want %f", result.Value, expectedAvg)
	}
}

func TestStudentAverageMultipleGrades(t *testing.T) {
	t.Parallel()

	student, err := gradebook.NewStudent("Michael", "Frede")
	if err != nil {
		t.Fatalf("NewStudent() returned error: %v", err)
	}

	student.GradesByCategory = map[string][]float64{
		"major": {88, 92, 85, 95},
	}

	result := student.Average("major")
	if !result.Valid {
		t.Error("Average(major) with grades should return Valid=true")
	}

	expectedAvg := 90.0
	if !floatEqual(result.Value, expectedAvg, 0.001) {
		t.Errorf("Average(major) with multiple grades = %f; want %f", result.Value, expectedAvg)
	}
}

func floatPtr(n float64) *float64 {
	return &n
}

func floatEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
