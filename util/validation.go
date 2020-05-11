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
	MsgURLPlaceholdersNotFound       = "Given URL '%s' doesn't contain placeholders. Templating will be ignored"
	MsgTemplateNotFound              = "Provided template file '%s' can not be found"
)

// Validator interface responsible for performing initial flags and templates validation
type Validator interface {

	// Validate method performs initial flags and templates validation
	Validate() *ValidationResult
}

// ValidatorProcessor interface responsible for processing validation errors and warnings
type ValidatorProcessor interface {

	// ProcessErrors method processes validation errors and warnings
	ProcessErrors()
}

// ValidationResult holds the information about validity of ValidatorEntity
type ValidationResult struct {
	valid        bool
	errMessages  []string
	warnMessages []string

	conf *app.Configuration
}

// ValidatorEntity holds the input flags and templates data
type ValidatorEntity struct {
	threads      int
	requestCount int
	sleep        int
	maxDuration  int
	url          string
	template     string

	conf *app.Configuration
}

// ProcessErrors method processes validation errors and warnings for ValidationResult
func (result *ValidationResult) ProcessErrors() {
	if len(result.errMessages) > 0 {
		_, _ = redColor.Println("The following validation errors occurred")

		for _, s := range result.errMessages {
			fmt.Println(redColor.Sprintf("   - %s", s))
		}
		fmt.Println()
	}

	if len(result.warnMessages) > 0 {
		_, _ = yellowColor.Println("The following validation warnings might impact an execution accuracy")

		for _, s := range result.warnMessages {
			fmt.Println(yellowColor.Sprintf("   - %s", s))
		}
		fmt.Println()
	}

	if !result.valid {
		os.Exit(1)
	}
}
