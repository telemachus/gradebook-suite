// Package gradebook is a library to read and write gradebook files.
package gradebook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	gradebookSuffix = ".gradebook"
	dateFmtLen      = len("YYYYMMDD")
)

// Term represents a grading period (e.g., a quarter or semester).
type Term struct {
	Start string
	End   string
}

// TermsByID maps Terms by short ID (e.g., "q1" points to the first quarter).
type TermsByID map[string]*Term

// AssignmentCategories stores basic assignment categories, the broad groupings
// such as "major" and "minor".
type AssignmentCategories []string

// LabelsByAssignmentCategory maps human-readable labels by a category of
// assignment. E.g., the category "cp" has the label "Class Participation",
// and the category "major" has the label "Major Assessments".
type LabelsByAssignmentCategory map[string]string

// WeightsByAssignmentCategory maps percentage values in a grading rubric by
// assignment category. The sum of the weights must equal 100 in order for this
// type to be valid.
type WeightsByAssignmentCategory map[string]int

// CategoriesByAssignmentType maps categories by their assignment type. (E.g.,
// "test", "essay", and "project" all have the category "major". Every
// assignment type must belong to one and only one category, and every
// assignment type must be present in CategoriesByAssignmentType. Also, and
// this is less obvious, every assignment category must have an assignment
// type. If a category has only a single assignment type, the category and type
// will often have the same name.  E.g., both the category and the type are
// "cp"."
type CategoriesByAssignmentType map[string]string

// Grade represents a single grade
type Grade struct {
	Email string   `json:"email"`
	Score *float64 `json:"score"`
}

// Grades stores Grade structs.
type Grades []*Grade

// Gradebook represents a single gradebook file.
type Gradebook struct {
	AssignmentDate     string `json:"assignment_date"`
	AssignmentName     string `json:"assignment_name"`
	AssignmentType     string `json:"assignment_type"`
	AssignmentCategory string `json:"assignment_category"`
	Grades             `json:"assignment_grades"`
}

// Student represents a student.
type Student struct {
	GradesByCategory map[string][]float64
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
}

// StudentsByEmail maps students by their email. (NB: an email is an
// appropriate equivalent to a database's primary key because emails are
// unique.)
type StudentsByEmail map[string]*Student

// Class represents a class and its students.
type Class struct {
	TermsByID                   `json:"terms_by_id"`
	LabelsByAssignmentCategory  `json:"labels_by_assignment_category"`
	WeightsByAssignmentCategory `json:"weights_by_assignment_category"`
	CategoriesByAssignmentType  `json:"categories_by_assignment_type"`
	StudentsByEmail             `json:"students_by_email"`
	Name                        string `json:"name"`
	AssignmentCategories        `json:"assignment_categories"`
}

// UnmarshalClass unmarshals a class.json file into a pointer to Class.
func UnmarshalClass(classFile string) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, err
	}

	var class Class
	err = json.Unmarshal(data, &class)
	if err != nil {
		return nil, err
	}

	return &class, nil
}

// UnmarshalGradebook unmarshals a gradebook file into a pointer to Gradebook.
func UnmarshalGradebook(gradebookFile string) (*Gradebook, error) {
	data, err := os.ReadFile(filepath.Clean(gradebookFile))
	if err != nil {
		return nil, err
	}

	var gradebook Gradebook
	err = json.Unmarshal(data, &gradebook)
	if err != nil {
		return nil, err
	}

	return &gradebook, nil
}

// dateSnip gets the date string from the end of a gradebook filename. If the
// function does not find a valid date (in YYYYMMDD format), then it return an
// error.
func dateSnip(dateStr string) (string, error) {
	dateStr = strings.TrimSuffix(dateStr, gradebookSuffix)
	dateStrLen := utf8.RuneCountInString(dateStr)
	if dateStrLen < dateFmtLen {
		return "", fmt.Errorf("[%s] does not contain a valid YYYYMMDD date", dateStr)
	}

	dateStr = dateStr[dateStrLen-dateFmtLen:]
	if _, err := time.Parse("20060102", dateStr); err != nil {
		return "", fmt.Errorf("[%s] does not contain a valid YYYYMMDD date", dateStr)
	}

	return dateStr, nil
}

// LoadGrades scans a given directory for *.gradebook files and adds grades
// from those files to students. The method returns an error if there is
// a problem reading, unmarshaling, or closing a file.
func (c *Class) LoadGrades(dir string, term *Term) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %q does not exist", dir)
	}

	gradebooks, err := filepath.Glob(filepath.Join(dir, "*.gradebook"))
	if err != nil {
		return err
	}

	for _, gradebook := range gradebooks {
		if term != nil {
			dateStr, err := dateSnip(gradebook)
			if err != nil {
				return err
			}

			if !term.Includes(dateStr) {
				continue
			}
		}

		if err := c.loadGradebookFile(gradebook); err != nil {
			return err
		}
	}

	return nil
}

func (c *Class) loadGradebookFile(gradebookPath string) error {
	gbData, err := UnmarshalGradebook(gradebookPath)
	if err != nil {
		return err
	}

	assignmentType := gbData.AssignmentType
	for _, grade := range gbData.Grades {
		if grade.Score == nil {
			continue
		}

		student, ok := c.StudentsByEmail[grade.Email]
		if !ok {
			return fmt.Errorf("no student with email %q", grade.Email)
		}

		category, ok := c.CategoriesByAssignmentType[assignmentType]
		if !ok {
			return fmt.Errorf("unrecognized assignment type %q", assignmentType)
		}
		_, ok = student.GradesByCategory[category]
		if !ok {
			return fmt.Errorf("unrecognized assignment category %q", category)
		}

		student.GradesByCategory[category] = append(student.GradesByCategory[category], *grade.Score)
	}

	return nil
}
