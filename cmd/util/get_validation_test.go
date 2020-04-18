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
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeThreadsWithOKOtherFlags(t *testing.T) {
	var givenZeroThreads = -1
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroThreadsWithOKOtherFlags(t *testing.T) {
	var givenZeroThreads = 0
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenZeroRequestCount = 0
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenZeroRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenZeroRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeRequestCountWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenNegativeRequestCount = -1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenNegativeRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenNegativeRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveSleepDelayWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenPositiveDelay = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenPositiveDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroSleepDelayWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenZeroDelay = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenZeroDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeSleepDelayWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenNegativeDelay = -1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenNegativeDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Delay in millis property", givenNegativeDelay) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveMaxDurationWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroMaxDurationWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.messages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeMaxDurationWithOKOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = -12

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.messages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Maximum execution duration property", givenMaxDuration) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}
