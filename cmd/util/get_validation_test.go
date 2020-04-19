package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePositiveThreadsWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeThreadsWithOkOtherFlags(t *testing.T) {
	var givenZeroThreads = -1
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroThreadsWithOkOtherFlags(t *testing.T) {
	var givenZeroThreads = 0
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveRequestCountWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroRequestCountWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenZeroRequestCount = 0
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenZeroRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenZeroRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeRequestCountWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenNegativeRequestCount = -1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenNegativeRequestCount).
		AddThreads(givenPositiveThreads).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenNegativeRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenPositiveDelay = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenPositiveDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenZeroDelay = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenZeroDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenNegativeDelay = -1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenNegativeDelay).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Delay in millis property", givenNegativeDelay) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = -12

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Maximum execution duration property", givenMaxDuration) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateEmptyTemplateAndUrlWithPlaceholders_WithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 666
	var givenEmptyTemplate = ""
	var givenUrlWithPlaceholders = "http://localhost:8080/get/#TE{1}/test/123"

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrlWithPlaceholders).
		AddTemplate(givenEmptyTemplate).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgUrlAddressInvalidWithReason, givenUrlWithPlaceholders, "URL address contains placeholder(s) for missing template file") {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNonExistingTemplateAndUrlWithPlaceholders_WithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 666
	var givenNotFoundTemplate = "test_NOT_FOUND.file"
	var givenUrlWithPlaceholders = "http://localhost:8080/get/#TE{1}/test/123"

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrlWithPlaceholders).
		AddTemplate(givenNotFoundTemplate).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgTemplateNotFound, givenNotFoundTemplate) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateExistingTemplateAndUrlWithPlaceholders_WithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 666
	var givenFoundTemplate = "test.file"
	var givenUrlWithPlaceholders = "http://localhost:1111/get/#TE{1}/test/123"

	// setup the file
	var testFileAbsPath, _ = filepath.Abs(givenFoundTemplate)
	var testFile, _ = os.Create(testFileAbsPath)
	_ = testFile.Close()

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrlWithPlaceholders).
		AddTemplate(givenFoundTemplate).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	_ = os.Remove(testFileAbsPath)
	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}
