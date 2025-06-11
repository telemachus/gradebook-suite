package gradebook

import (
	"fmt"
	"strconv"
)

// AverageResult represents the result of *Student.Average. The result is only
// valid for use if the given Student has scores of a given type to average.
type AverageResult struct {
	Value float64
	Valid bool
}

// String returns a string representation of an AverageResult.
func (ar AverageResult) String() string {
	if ar.Valid {
		return strconv.FormatFloat(ar.Value, 'g', -1, 64)
	}

	return "No results"
}

// Average calculates the average for a slice of scores in a given category and
// returns an AverageResult and an error. If the category is unknown, the
// method returns an error. If the slice of scores is empty, the method returns
// an invalid AverageResult.
func (s *Student) Average(category string) (AverageResult, error) {
	if _, ok := s.GradesByCategory[category]; !ok {
		return AverageResult{Valid: false}, fmt.Errorf("unknown grade type: %q", category)
	}

	if len(s.GradesByCategory[category]) < 1 {
		return AverageResult{Valid: false}, nil
	}

	return AverageResult{Value: fmean(s.GradesByCategory[category]), Valid: true}, nil
}

// TotalAverage returns an AverageResult and error for all of a student's
// scores. The method will return an invalid result if the student has no
// scores in any category. The method will return an error if any call to
// Average for a given category returns an error.
func (s *Student) TotalAverage(weights WeightsByAssignmentCategory) (AverageResult, error) {
	var summedAverage float64
	var summedWeight int

	for assignmentType, weight := range weights {
		typeAverage, err := s.Average(assignmentType)
		if err != nil {
			return AverageResult{Valid: false}, err
		}

		if typeAverage.Valid {
			summedAverage += typeAverage.Value * float64(weight)
			summedWeight += weight
		}
	}

	if summedWeight == 0 {
		return AverageResult{Valid: false}, nil
	}

	totalAverage := summedAverage / float64(summedWeight)

	return AverageResult{Value: totalAverage, Valid: true}, nil
}

// Fmean will panic if scores is empty.
func fmean(scores []float64) float64 {
	var sum float64
	for _, score := range scores {
		sum += score
	}

	return sum / float64(len(scores))
}
