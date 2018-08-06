package params

// Context ...
type Context interface {
	Set(string, string) error

	Bool(string) bool
	Int(string) int
	Int64(string) int64
	String(string) string
	StringSlice(string) []string
}
