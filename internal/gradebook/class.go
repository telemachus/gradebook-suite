package gradebook

import (
	"cmp"
	"maps"
	"slices"
)

// EmailsSortedByStudentName returns a slice of student emails sorted by student name.
func (c *Class) EmailsSortedByStudentName() []string {
	if len(c.StudentsByEmail) == 0 {
		return []string{}
	}

	emails := slices.Collect(maps.Keys(c.StudentsByEmail))
	slices.SortFunc(emails, func(emailA, emailB string) int {
		studentA := c.StudentsByEmail[emailA]
		studentB := c.StudentsByEmail[emailB]
		return cmpStudent(studentA, studentB)
	})

	return emails
}

// StudentsSortedByName returns a slice of students sorted by last and first name.
func (c *Class) StudentsSortedByName() []*Student {
	students := make([]*Student, 0, len(c.StudentsByEmail))
	for _, student := range c.StudentsByEmail {
		students = append(students, student)
	}
	slices.SortFunc(students, cmpStudent)
	return students
}

func cmpStudent(studentA, studentB *Student) int {
	if studentA.LastName == studentB.LastName {
		return cmp.Compare(studentA.FirstName, studentB.FirstName)
	}

	return cmp.Compare(studentA.LastName, studentB.LastName)
}
