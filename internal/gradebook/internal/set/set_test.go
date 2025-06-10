package set_test

import (
	"testing"

	"github.com/telemachus/gradebook-suite/internal/gradebook/internal/set"
)

func TestSetEquals(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		lhs  set.Set[int]
		rhs  set.Set[int]
		want bool
	}{
		"empty sets are equal": {
			lhs:  set.New[int](),
			rhs:  set.New[int](),
			want: true,
		},
		"equal one-item sets are equal": {
			lhs:  set.New(1),
			rhs:  set.New(1),
			want: true,
		},
		"equal multi-item sets are equal": {
			lhs:  set.New(1, 2, 3),
			rhs:  set.New(1, 2, 3),
			want: true,
		},
		"equal sets are equal regardless of declaration duplicates": {
			lhs:  set.New(1, 1, 1, 2, 4),
			rhs:  set.New(2, 2, 2, 4, 4, 1),
			want: true,
		},
		"empty set is unequal to set with elements": {
			lhs:  set.New[int](),
			rhs:  set.New(1, 2),
			want: false,
		},
		"unequal sets are unequal": {
			lhs:  set.New(1, 2, 4),
			rhs:  set.New(1, 2),
			want: false,
		},
	}

	for msg, tt := range tests {
		tt := tt

		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tt.lhs.Equals(tt.rhs)
			if got != tt.want {
				t.Errorf("%s.equals(%s) = %t; want %t", tt.lhs, tt.rhs, got, tt.want)
			}
		})
	}
}

func TestSetString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		set  set.Set[int]
		want string
	}{
		"empty set.String == {}": {
			set:  set.New[int](),
			want: "{}",
		},
		"set(1).String = {1}": {
			set:  set.New(1),
			want: "{1}",
		},
	}

	for msg, tt := range tests {
		tt := tt

		t.Run(msg, func(t *testing.T) {
			t.Parallel()

			got := tt.set.String()
			if got != tt.want {
				t.Errorf("set.String() = %q; want %q", got, tt.want)
			}
		})
	}
}
