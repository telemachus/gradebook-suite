package gradebook

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UnmarshalCalcClass unmarshals a class.json file into a pointer to Class.
// Unlike UnmarshalClass, this function creates the Grades map needed to store
// grades.
func UnmarshalCalcClass(classFile string) (*Class, error) {
	data, err := os.ReadFile(filepath.Clean(classFile))
	if err != nil {
		return nil, err
	}

	var class Class
	err = json.Unmarshal(data, &class)
	if err != nil {
		return nil, err
	}

	for name, student := range class.StudentsByEmail {
		gradesByType := make(map[string][]float64, len(class.AssignmentTypes))
		for _, cat := range class.AssignmentTypes {
			// TODO: how should I decide the capacity here?
			gradesByType[cat] = make([]float64, 0, 25)
		}

		student.GradesByType = gradesByType
		class.StudentsByEmail[name] = student
	}

	return &class, nil
}
