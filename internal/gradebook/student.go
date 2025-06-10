package gradebook

import (
	"errors"
)

// NewStudent a new *Student. If firstName or lastName is empty, the method
// returns an error.
func NewStudent(firstName, lastName string) (*Student, error) {
	if firstName == "" {
		return nil, errors.New("firstName cannot be an empty string")
	} else if lastName == "" {
		return nil, errors.New("lastName cannot be an empty string")
	}

	return &Student{FirstName: firstName, LastName: lastName}, nil
}
