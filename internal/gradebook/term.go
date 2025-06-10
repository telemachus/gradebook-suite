package gradebook

// Includes determines whether a date falls within a term.
func (t Term) Includes(d string) bool {
	if d < t.Start || d > t.End {
		return false
	}

	return true
}
