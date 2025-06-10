package gradebook

import (
	"testing"
)

func TestDateSnipValid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileName string
		want     string
	}{
		"hw-aeneid-ii-200-220-20230907.gradebook": {
			fileName: "hw-aeneid-ii-200-220-20230907.gradebook",
			want:     "20230907",
		},
		"quiz-aeneid-ii-200-220-20241027.gradebook": {
			fileName: "quiz-aeneid-ii-200-220-20241027.gradebook",
			want:     "20241027",
		},
		"20230213.gradebook": {
			fileName: "20230213.gradebook",
			want:     "20230213",
		},
	}

	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got, err := dateSnip(tc.fileName)
			if err != nil {
				t.Fatalf("err should be nil; got %s", err)
			}

			if got != tc.want {
				t.Errorf("dateSnip(%s) returns %s; want %s", tc.fileName, got, tc.want)
			}
		})
	}
}

func TestDateSnipInvalid(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		fileName string
	}{
		"invalid date (correct length and placement)": {
			fileName: "hw-aeneid-ii-200-220-00000000.gradebook",
		},
		"seven digit date": {
			fileName: "2024102.gradebook",
		},
		"[empty string]": {
			fileName: "",
		},
	}
	for msg, tc := range testCases {
		t.Run(msg, func(t *testing.T) {
			t.Parallel()
			_, err := dateSnip(tc.fileName)
			if err == nil {
				t.Errorf("err is nil; should have error for [%s]", tc.fileName)
			}
		})
	}
}
