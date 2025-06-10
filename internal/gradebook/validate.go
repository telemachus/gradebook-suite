package gradebook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/telemachus/gradebook-suite/internal/gradebook/internal/set"
	"golang.org/x/exp/maps"
)

func zvalErr(zvals []string) error {
	switch len(zvals) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("gradebook: a field in Class is unset: %s", zvals[0])
	default:
		return fmt.Errorf("gradebooks: fields in Class are unset: %s", strings.Join(zvals, ", "))
	}
}

// checkInitialization ensures that a *Class has no dangerous zero values.
func (c *Class) checkInitialization() error {
	zvals := make([]string, 0, 7)

	if c.Name == "" {
		zvals = append(zvals, "Name")
	}
	if c.TermsByID == nil {
		zvals = append(zvals, "TermsByID")
	}
	if c.AssignmentTypes == nil {
		zvals = append(zvals, "AssignmentTypes")
	}
	if c.LabelsByAssignmentType == nil {
		zvals = append(zvals, "LabelsByAssignmentType")
	}
	if c.WeightsByAssignmentType == nil {
		zvals = append(zvals, "WeightsByAssignmentType")
	}
	if c.CategoriesByAssignmentType == nil {
		zvals = append(zvals, "CategoriesByAssignmentType")
	}
	if c.StudentsByEmail == nil {
		zvals = append(zvals, "StudentsByEmail")
	}

	return zvalErr(zvals)
}

// checkWeightsSum ensures that c.Weights adds up to 100%.
func (c *Class) checkWeightsSum() error {
	total := 0
	for _, n := range c.WeightsByAssignmentType {
		total += n
	}

	if total != 100 {
		return errors.New("gradebook: WeightsByAssignmentType must equal 100%")
	}

	return nil
}

// checkEq returns an error if two sets are not equal or nil if they are.
func checkEq[T comparable](lhs, rhs set.Set[T]) error {
	if !lhs.Equals(rhs) {
		return fmt.Errorf("%s and %s are not equal sets", lhs, rhs)
	}

	return nil
}

// Validate checks whether a *Class is valid. It returns nil if the *Class is
// valid. Otherwise it returns an error containing one more errors from the
// individual checks. Those errors are combined using errors.Join.
func (c *Class) Validate() error {
	assignmentsSet := set.New(c.AssignmentTypes...)
	categoriesSet := set.New(maps.Values(c.CategoriesByAssignmentType)...)
	weightsSet := set.New(maps.Keys(c.WeightsByAssignmentType)...)
	labelsSet := set.New(maps.Keys(c.LabelsByAssignmentType)...)

	return errors.Join(
		c.checkInitialization(),
		c.checkWeightsSum(),
		checkEq(assignmentsSet, categoriesSet),
		checkEq(assignmentsSet, labelsSet),
		checkEq(assignmentsSet, weightsSet),
	)
}
