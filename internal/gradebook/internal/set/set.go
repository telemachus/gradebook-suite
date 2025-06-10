// Package set provides a minimal ADT and operations for sets.
package set

import (
	"fmt"
	"strings"
)

// Set represents a set of comparable items of type E.
type Set[E comparable] map[E]struct{}

// New returns a set of type E containing all items in elems.
func New[E comparable](elems ...E) Set[E] {
	s := make(Set[E], len(elems))
	for _, el := range elems {
		s[el] = struct{}{}
	}

	return s
}

// Equals determines whether s and other are equal sets.
func (s Set[E]) Equals(other Set[E]) bool {
	if len(s) != len(other) {
		return false
	}

	for el := range s {
		_, ok := other[el]
		if !ok {
			return false
		}
	}

	return true
}

// String returns a string representation of s.
func (s Set[E]) String() string {
	elems := make([]string, 0, len(s))
	for el := range s {
		elems = append(elems, fmt.Sprintf("%v", el))
	}

	return fmt.Sprintf("{%s}", strings.Join(elems, ", "))
}
