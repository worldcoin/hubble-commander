package utils

// Why is this even needed?
// https://github.com/ksimka/go-is-not-good
func MakeStringPointer(s string) *string {
	return &s
}
