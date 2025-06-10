package gradebook_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/gradebook-suite/internal/gradebook"
)

func TestUnmarshalCalcClass(t *testing.T) {
	t.Parallel()

	want := fakeCalcClass()
	got, err := gradebook.UnmarshalCalcClass(classJSON)
	if err != nil {
		t.Fatalf("gradebook.UnmarshalClass(classJSON) = %v; want nil error", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("gradebook.UnmarshalClass(classJSON) mismatch(-want +got):\n%s", diff)
	}
}

func fakeGradesMap() map[string][]float64 {
	gradesByCategories := make(map[string][]float64, 3)
	gradesByCategories["major"] = make([]float64, 0, 50)
	gradesByCategories["minor"] = make([]float64, 0, 50)
	gradesByCategories["cp"] = make([]float64, 0, 50)

	return gradesByCategories
}
