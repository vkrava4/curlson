package util

const (
	MsgShouldBePositive = "A number of %s should be positive. Currently it's: '%d'"
)

type Validator interface {
	Validate() *ValidationResult
}

type ValidationResult struct {
	valid    bool
	messages []string
}

type ValidatorEntity struct {
	threads      int
	requestCount int
	url          string
	template     string
}
