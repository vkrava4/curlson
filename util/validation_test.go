package util

import "testing"

func TestValidationResult_ProcessErrorsWhenErrorsAndWarnListEmpty(t *testing.T) {
	var result = &ValidationResult{
		valid:        true,
		errMessages:  []string{},
		warnMessages: []string{},
		conf:         nil,
	}

	result.ProcessErrors()
}

func TestValidationResult_ProcessErrorsWhenErrorsAndWarnListNil(t *testing.T) {
	var result = &ValidationResult{
		valid:        true,
		errMessages:  nil,
		warnMessages: nil,
		conf:         nil,
	}

	result.ProcessErrors()
}

func TestValidationResult_ProcessErrorsWhenErrorsAndWarnListNotEmpty(t *testing.T) {
	var result = &ValidationResult{
		valid:        true,
		errMessages:  []string{"ERROR"},
		warnMessages: []string{"WARN"},
		conf:         nil,
	}

	result.ProcessErrors()
}
