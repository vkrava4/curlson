package util

const (
	// Generic validation constants
	MsgShouldBePositive       = "%s should be positive. Currently it's: '%d'"
	MsgShouldBePositiveOrZero = "%s should be positive or equal to zero. Currently it's: '%d'"

	// URL-related validation constants
	MsgUrlAddressInvalid           = "Provided URL address is invalid. Currently it's: '%s'"
	MsgUrlAddressInvalidWithReason = "Provided URL address: '%s' is invalid. Reason: %s"

	// Template-related validation constants
	MsgTemplateInvalidWithReason     = "Provided template file '%s' is invalid. Reason: %s"
	MsgTemplatePathInvalidWithReason = "Provided template file path '%s' is invalid. Reason: %s"
	MsgCantOpenTemplateWithReason    = "Provided template file '%s' can not be opened. Reason: %s"
	MsgTemplateNotFound              = "Provided template file '%s' can not be found"
)

type Validator interface {
	Validate() *ValidationResult
}

type ValidationResult struct {
	valid        bool
	errMessages  []string
	warnMessages []string

	conf *AppConfiguration
}

type ValidatorEntity struct {
	threads      int
	requestCount int
	sleep        int
	maxDuration  int
	url          string
	template     string

	conf *AppConfiguration
}
