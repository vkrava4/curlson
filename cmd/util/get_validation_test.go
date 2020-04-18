package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidatePositiveThreadsWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).AddThreads(givenPositiveThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeThreadsWithOKOtherFlags(t *testing.T) {
	var givenZeroThreads = -1
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).AddThreads(givenZeroThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") != fmt.Sprintf(MsgShouldBePositive, "threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroThreadsWithOKOtherFlags(t *testing.T) {
	var givenZeroThreads = 0
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).AddThreads(givenZeroThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") != fmt.Sprintf(MsgShouldBePositive, "threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).AddThreads(givenPositiveThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenZeroRequestCount = 0
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenZeroRequestCount).AddThreads(givenPositiveThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") != fmt.Sprintf(MsgShouldBePositive, "requests per thread", givenZeroRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenNegativeRequestCount = -1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenNegativeRequestCount).AddThreads(givenPositiveThreads).Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") != fmt.Sprintf(MsgShouldBePositive, "requests per thread", givenNegativeRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}
