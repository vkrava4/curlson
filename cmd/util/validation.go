package util

const (
	MsgShouldBePositive       = "%s should be positive. Currently it's: '%d'"
	MsgShouldBePositiveOrZero = "%s should be positive or equal to zero. Currently it's: '%d'"
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
	sleep        int
	maxDuration  int
	url          string
	template     string
}
