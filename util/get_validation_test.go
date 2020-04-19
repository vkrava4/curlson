package util

import (
	"fmt"
	"github.com/vkrava4/curlson/app"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestValidatePositiveThreadsWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeThreadsWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenZeroThreads = -1
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroThreadsWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenZeroThreads = 0
	var givenPositiveRequestCount = 999
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenZeroThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of threads", givenZeroThreads) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveRequestCountWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroRequestCountWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenZeroRequestCount = 0
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenZeroRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenZeroRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeRequestCountWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenNegativeRequestCount = -1
	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenNegativeRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositive, "Amount of requests per thread", givenNegativeRequestCount) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenPositiveDelay = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenPositiveDelay).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenZeroDelay = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenZeroDelay).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeSleepDelayWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenNegativeDelay = -1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddSleep(givenNegativeDelay).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Delay in millis property", givenNegativeDelay) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidatePositiveMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 1

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateZeroMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = 0

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateNegativeMaxDurationWithOkOtherFlags(t *testing.T) {
	var givenUrl = "http://localhost:8080"
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 1
	var givenMaxDuration = -12

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddMaxDuration(givenMaxDuration).
		AddUrl(givenUrl).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgShouldBePositiveOrZero, "Maximum execution duration property", givenMaxDuration) {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}
}

func TestValidateEmptyTemplateAndWithEmptyUrl_WithOkOtherFlags(t *testing.T) {
	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 666
	var givenEmptyTemplate = ""
	var givenEmptyUrl = ""

	var getValidator = &GetValidator{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenEmptyUrl).
		AddTemplate(givenEmptyTemplate).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	if actualValidationResult.valid || strings.Join(actualValidationResult.errMessages, ",") !=
		fmt.Sprintf(MsgUrlAddressInvalidWithReason, givenEmptyUrl, fmt.Sprintf("A string: '%s' is not valid URL", givenEmptyUrl)) {
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
		fmt.Sprintf(MsgUrlAddressInvalidWithReason, givenUrlWithPlaceholders, fmt.Sprintf("Given URL: '%s' has unresolved placeholders", givenUrlWithPlaceholders)) {
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
	var givenUrlWithPlaceholders = "http://localhost:1111/get/#T{0}/test/123?q=#T{1}"

	// setup the file
	var testFileAbsPath, _ = filepath.Abs(givenFoundTemplate)

	_ = ioutil.WriteFile(testFileAbsPath, []byte("TEST,ONE\n"), filesMode)

	var getValidator = &GetValidator{}
	var appConf = &app.Configuration{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrlWithPlaceholders).
		AddTemplate(givenFoundTemplate).
		WithAppConfiguration(appConf).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	_ = os.Remove(testFileAbsPath)
	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}

	if !appConf.Template.Enabled || appConf.Template.Size != 1 || appConf.Template.Path != testFileAbsPath {
		t.Errorf("Unexpected app configuration result %v", actualValidationResult)
	}
}

func TestValidateExistingTemplateAndUrlWithPlaceholders_WithOkOtherFlags_ForLargeDataSet(t *testing.T) {
	var givenNumberOfRecords = 10000

	var givenPositiveThreads = 999
	var givenPositiveRequestCount = 666
	var givenFoundTemplate = "test.file"
	var givenUrlWithPlaceholders = "http://localhost:1111/get/#T{6}/test/123?q=#T{1}"

	// setup the file
	var testFileAbsPath, _ = filepath.Abs(givenFoundTemplate)
	var testFile, _ = os.Create(testFileAbsPath)

	for i := 0; i < givenNumberOfRecords; i++ {
		_, _ = testFile.WriteString(fmt.Sprintf("TEST%d,ONE%d,THREE,FOUR,test123,777,VLADKRAVA\n", i, i))
	}

	var start = time.Now()
	var getValidator = &GetValidator{}
	var appConf = &app.Configuration{}
	var validatorEntity = getValidator.AddRequestCount(givenPositiveRequestCount).
		AddThreads(givenPositiveThreads).
		AddUrl(givenUrlWithPlaceholders).
		AddTemplate(givenFoundTemplate).
		WithAppConfiguration(appConf).
		Entity()

	var actualValidationResult = validatorEntity.Validate()

	_ = os.Remove(testFileAbsPath)
	if !actualValidationResult.valid || len(actualValidationResult.errMessages) > 0 {
		t.Errorf("Unexpected validation result %v", actualValidationResult)
	}

	if !appConf.Template.Enabled || appConf.Template.Size != givenNumberOfRecords || appConf.Template.Path != testFileAbsPath {
		t.Errorf("Unexpected app configuration result %v", actualValidationResult)
	}

	t.Logf("For %d items template items validation took %d ms", givenNumberOfRecords, time.Now().Sub(start).Milliseconds())
}
