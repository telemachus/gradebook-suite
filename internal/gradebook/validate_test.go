package gradebook_test

import (
	"testing"

	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

func TestValid(t *testing.T) {
	t.Parallel()

	class := fakeClass()
	if err := class.Validate(); err != nil {
		t.Errorf("class.ZeroValues() = %v; want no error", err)
	}
}

func TestInitializationInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"class.Name unset": {
			transformClass: func(c *gradebook.Class) {
				c.Name = ""
			},
		},
		"class.TermsByID unset": {
			transformClass: func(c *gradebook.Class) {
				c.TermsByID = nil
			},
		},
		"class.AssignmentTypes unset": {
			transformClass: func(c *gradebook.Class) {
				c.AssignmentTypes = nil
			},
		},
		"class.LabelsByAssignmentType unset": {
			transformClass: func(c *gradebook.Class) {
				c.LabelsByAssignmentType = nil
			},
		},
		"class.WeightsByAssignmentType unset": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType = nil
			},
		},
		"class.CategoriesByAssignmentType unset": {
			transformClass: func(c *gradebook.Class) {
				c.CategoriesByAssignmentType = nil
			},
		},
		"class.StudentsByEmail unset": {
			transformClass: func(c *gradebook.Class) {
				c.StudentsByEmail = nil
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Error("c.Validate() returns nil; want error for zero value(s)")
			}
		})
	}
}

func TestWeightsSumInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"no WeightsByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType = gradebook.WeightsByAssignmentType{}
			},
		},
		"WeightsByAssignmentType under 100": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType["major"] = 25
			},
		},
		"WeightsByAssignmentType over 100": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType["major"] = 75
			},
		},
		"Weights below 0": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType["major"] = -175
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Error("class.Validate() returns nil; want error for incorrect Weights")
			}
		})
	}
}

func TestSetEqualityInvalid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		transformClass func(c *gradebook.Class)
	}{
		"missing items from AssignmentTypes": {
			transformClass: func(c *gradebook.Class) {
				clear(c.AssignmentTypes)
			},
		},
		"extra item in AssignmentTypes": {
			transformClass: func(c *gradebook.Class) {
				c.AssignmentTypes = append(c.AssignmentTypes, "random")
			},
		},
		"missing item from CategoriesByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				delete(c.CategoriesByAssignmentType, "cp")
			},
		},
		"extra item in CategoriesByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				c.CategoriesByAssignmentType["random"] = "random"
			},
		},
		"missing item from LabelsByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				delete(c.LabelsByAssignmentType, "cp")
			},
		},
		"extra item in LabelsByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				c.LabelsByAssignmentType["random"] = "Random Item"
			},
		},
		"missing item from WeightsByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				delete(c.WeightsByAssignmentType, c.AssignmentTypes[0])
			},
		},
		"extra item in WeightsByAssignmentType": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentType["random"] = 0
			},
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			class := fakeClass()
			if tc.transformClass != nil {
				tc.transformClass(class)
			}

			if err := class.Validate(); err == nil {
				t.Errorf("class.Validate() returns nil; want error for %s", msg)
			}
		})
	}
}
