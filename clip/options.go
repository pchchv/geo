package clip

// Option is a possible parameter to the clip operations.
type Option func(*options)

type options struct {
	openBound bool
}
