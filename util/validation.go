package util

import (
	"fmt"
	"github.com/vkrava4/curlson/app"
	"os"
)

const (
	// Generic validation constants
	MsgShouldBePositive       = "%s should be positive. Currently it's: '%d'"
	MsgShouldBePositiveOrZero = "%s should be positive or equal to zero. Currently it's: '%d'"

	// URL-related validation constants
	MsgURLAddressInvalidWithReason = "Provided URL address: '%s' is invalid. Reason: %s"

	// Template-related validation constants
	MsgTemplatePathInvalidWithReason = "Provided template file path '%s' is invalid. Reason: %s"
	MsgCantOpenTemplateWithReason    = "Provided template file '%s' can not be opened. Reason: %s"
	MsgTemplateNotFound              = "Provided template file '%s' can not be found"
)

type Validator interface {
	Validate() *ValidationResult
}

type ValidatorProcessor interface {
	ProcessErrors()
}

type ValidationResult struct {
	valid        bool
	errMessages  []string
	warnMessages []string

	conf *app.Configuration
}

type ValidatorEntity struct {
	threads      int
	requestCount int
	sleep        int
	maxDuration  int
	url          string
	template     string

	conf *app.Configuration
}

func (result *ValidationResult) ProcessErrors() {
	if len(result.warnMessages) > 0 {
		fmt.Println()
		_, _ = yellowColor.Println("The following validation warnings occurred")

		for _, s := range result.warnMessages {
			fmt.Println(yellowColor.Sprintf("   - %s", s))
		}
	}

	if len(result.errMessages) > 0 {
		fmt.Println()
		_, _ = redColor.Println("The following validation errors occurred")

		for _, s := range result.errMessages {
			fmt.Println(redColor.Sprintf("   - %s", s))
		}
	}

	if !result.valid {
		os.Exit(1)
	}
}
