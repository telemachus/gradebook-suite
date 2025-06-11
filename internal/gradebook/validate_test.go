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
		"class.AssignmentCategories unset": {
			transformClass: func(c *gradebook.Class) {
				c.AssignmentCategories = nil
			},
		},
		"class.LabelsByAssignmentCategory unset": {
			transformClass: func(c *gradebook.Class) {
				c.LabelsByAssignmentCategory = nil
			},
		},
		"class.WeightsByAssignmentCategory unset": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory = nil
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
		"no WeightsByAssignmentCategory": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory = gradebook.WeightsByAssignmentCategory{}
			},
		},
		"WeightsByAssignmentCategory under 100": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory["major"] = 25
			},
		},
		"WeightsByAssignmentCategory over 100": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory["major"] = 75
			},
		},
		"Weights below 0": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory["major"] = -175
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
		"missing items from AssignmentCategories": {
			transformClass: func(c *gradebook.Class) {
				clear(c.AssignmentCategories)
			},
		},
		"extra item in AssignmentCategories": {
			transformClass: func(c *gradebook.Class) {
				c.AssignmentCategories = append(c.AssignmentCategories, "random")
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
		"missing item from LabelsByAssignmentCategory": {
			transformClass: func(c *gradebook.Class) {
				delete(c.LabelsByAssignmentCategory, "cp")
			},
		},
		"extra item in LabelsByAssignmentCategory": {
			transformClass: func(c *gradebook.Class) {
				c.LabelsByAssignmentCategory["random"] = "Random Item"
			},
		},
		"missing item from WeightsByAssignmentCategory": {
			transformClass: func(c *gradebook.Class) {
				delete(c.WeightsByAssignmentCategory, c.AssignmentCategories[0])
			},
		},
		"extra item in WeightsByAssignmentCategory": {
			transformClass: func(c *gradebook.Class) {
				c.WeightsByAssignmentCategory["random"] = 0
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
