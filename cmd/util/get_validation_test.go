package util

import (
	"testing"
)

func TestGetValidatorForPositiveThreadsWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).AddThreads(givenPositiveThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}
