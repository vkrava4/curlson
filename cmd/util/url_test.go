package util

import (
	"fmt"
	"testing"
)

func TestPrepareUrlWithRegularAttributes(t *testing.T) {
	var givenUrlTemplate = "http://localhost:8080/get/echo/with_status_param_and_request_param/#T{0}?data=#T{1}"
	var givenValueLine = "Test,Test2"

	var expectedUrl = "http://localhost:8080/get/echo/with_status_param_and_request_param/Test?data=Test2"
	var actualUrl, _ = PrepareUrl(givenUrlTemplate, givenValueLine)

	if expectedUrl != actualUrl {
		t.Errorf("PrepareUrl result is incorrect, actual: '%s', expected: '%s'", actualUrl, expectedUrl)
	}
}

func TestNotPrepareUrlWithRegularAttributesIfContainsTemplatePlaceholders(t *testing.T) {
	var givenUrlTemplate = "http://localhost:8080/get/echo/with_status_param_and_request_param/#T{0}/#T{0}/#T{2}?data=#T{1}"
	var givenValueLine = "Test,Test2"

	var expectedUrlInError = "http://localhost:8080/get/echo/with_status_param_and_request_param/Test/Test/#T{2}?data=Test2"
	var _, err = PrepareUrl(givenUrlTemplate, givenValueLine)

	if err != nil && err.Error() == fmt.Sprintf("Giving URL: '%s' has unresolved placeholders", expectedUrlInError) {
	} else {
		t.Error("PrepareUrl result is incorrect for this test case, it should have an error")
	}
}

func TestPrepareNonUrlWithRegularAndEscapedAttributes(t *testing.T) {
	var givenUrlTemplate = "htp://localhost:8080/get/echo/with_status_param_and_request_param/#T{0}?data=#TE{1}"
	var givenValueLine = "test/200/123,T  H//I // S\\\\I \\ S++T + E   ST1"

	var expectedErr = "A string: 'htp://localhost:8080/get/echo/with_status_param_and_request_param/test/200/123?data=T++H%2F%2FI+%2F%2F+S%5C%5CI+%5C+S%2B%2BT+%2B+E+++ST1' is not valid URL"
	var _, err = PrepareUrl(givenUrlTemplate, givenValueLine)

	if err != nil && err.Error() == expectedErr {
	} else {
		t.Error("PrepareUrl result is incorrect for this test case, it should have an error")
	}
}

func TestPrepareUrlWithRegularAndEscapedAttributes(t *testing.T) {
	var givenUrlTemplate = "http://localhost:8080/get/echo/with_status_param_and_request_param/#T{0}?data=#TE{1}"
	var givenValueLine = "test/200/123,T  H//I // S\\\\I \\ S++T + E   ST1"

	var expectedUrl = "http://localhost:8080/get/echo/with_status_param_and_request_param/test/200/123?data=T++H%2F%2FI+%2F%2F+S%5C%5CI+%5C+S%2B%2BT+%2B+E+++ST1"
	var actualUrl, _ = PrepareUrl(givenUrlTemplate, givenValueLine)

	if expectedUrl != actualUrl {
		t.Errorf("PrepareUrl result is incorrect, actual: '%s', expected: '%s'", actualUrl, expectedUrl)
	}
}

func TestContainsTemplatePlaceholders(t *testing.T) {
	var urlWithTemplatePlaceholder = "http://localhost:8080/get/#T{100}/test/123"
	var urlWithTemplateEscapedPlaceholder = "http://localhost:8080/get/#TE{1}/test/123"
	var urlWithBothPlaceholders = "http://localhost:8080/get/#TE{1}/test/123?q#T{0}"
	var urlWithNonePlaceholders = "http://localhost:8080/get/#TE/test/123?q1"

	if !ContainsTemplatePlaceholders(urlWithTemplatePlaceholder) || !ContainsTemplatePlaceholders(urlWithTemplateEscapedPlaceholder) ||
		!ContainsTemplatePlaceholders(urlWithBothPlaceholders) || ContainsTemplatePlaceholders(urlWithNonePlaceholders) {
		t.Error("Can not successfully determine whether Url contains Template placeholders or not")
	}
}

func TestParseAndValidateUrlForValidStrings(t *testing.T) {
	var validUrls = []string{
		"http://localhost",
		"https://localhost:9999",
		"http://localhost/path/path/path",
		"http://localhost/path?data=1&data2=2",
	}

	for i, url := range validUrls {
		var parsedUrl, _ = ParseAndValidateUrl(url)
		if parsedUrl != validUrls[i] {
			t.Errorf("ParseAndValidateUrl result is incorrect, actual: '%s', expected: '%s'", parsedUrl, validUrls[i])
		}
	}
}

func TestParseAndValidateUrlForInvalidStrings(t *testing.T) {
	var validUrls = []string{
		"ttp://localhost",
		"httpss://localhost:9999",
		"://localhost/path/path/path",
		"localhost/path?data=1&data2=2",
		"http:// /path?data=1&data2=2",
	}

	for i, url := range validUrls {
		var _, err = ParseAndValidateUrl(url)

		if err != nil && err.Error() == fmt.Sprintf("A string: '%s' is not valid URL", validUrls[i]) {
		} else {
			t.Error("ParseAndValidateUrl result is incorrect for this test case, it should have an error")
		}
	}
}
