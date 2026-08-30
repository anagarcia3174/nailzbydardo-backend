package testutil

import "time"

// imports: "testing", "time"

// strPtr returns a pointer to the given string — a small helper since
// Go doesn't allow taking the address of a string literal directly.
func StrPtr(s string) *string { return &s }


// timePtr returns a pointer to the given time.Time, same reasoning.
func TimePtr(t time.Time) *time.Time { return &t }


// intPtr / int64Ptr — you'll likely want one of these soon too, given
// how many of your model fields (LateFee, Tip, DesignPrice references,
// etc.) are *int64. Worth adding proactively:
func Int64Ptr(i int64) *int64 { return &i }


// EqualStrPtr compares two *string values for equality, treating nil
// as only equal to nil (not to an empty string) — the correct semantics
// for your nullable fields.
func EqualStrPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}


// EqualInt64Ptr — same pattern, for your *int64 fields (Tip, LateFee, etc.)
func EqualInt64Ptr(a, b *int64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}


// EqualTimePtr compares two *time.Time values, using time.Time's own
// .Equal method (NOT ==, since two time.Time values representing the
// same instant can still fail a == comparison due to internal
// monotonic-clock/location differences — .Equal is the correct way
// to compare times in Go).
func EqualTimePtr(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}


// DerefStr returns the dereferenced string, or a placeholder like
// "<nil>" if the pointer is nil — useful for producing readable
// t.Errorf messages without a nil-check at every call site.
func DerefStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}