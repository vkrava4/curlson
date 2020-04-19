package util

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	templateEscapedRegex   = "(\\#TE{\\d+})"
	templateUnescapedRegex = "(\\#T{\\d+})"
)

// Prepares given url param (URL address) by replacing placeholder values: `#TE{i}` or `#T{i}` with given `valuesLine[i]` value
// #TE{i} - Template Escaped, means that its values will be escaped so it can be safely placed inside a URL path segment. See: url.QueryEscape(value)
// #T{i} - Template Unescaped, means that its values will be placed to resulted URL in a raw format
func PrepareUrl(urlTemplate string, valuesLine string) (string, error) {
	var values = strings.Split(valuesLine, ",")
	var valuesLen = len(values)

	if valuesLen != 0 {
		for i, value := range values {
			urlTemplate = strings.ReplaceAll(urlTemplate, fmt.Sprintf("#TE{%d}", i), url.QueryEscape(value))
			urlTemplate = strings.ReplaceAll(urlTemplate, fmt.Sprintf("#T{%d}", i), value)
		}
	}

	if ContainsTemplatePlaceholders(urlTemplate) {
		return "", errors.New(fmt.Sprintf("Given URL: '%s' has unresolved placeholders", urlTemplate))
	}

	return ParseAndValidateUrl(urlTemplate)
}

func ParseAndValidateUrl(urlAddress string) (string, error) {
	var parsedUrl, err = url.ParseRequestURI(urlAddress)
	if err == nil && (parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https") && parsedUrl.Host != "" {
		return parsedUrl.String(), nil
	}

	return "", errors.New(fmt.Sprintf("A string: '%s' is not valid URL", urlAddress))
}

func ContainsTemplatePlaceholders(urlAddress string) bool {
	return containsTemplatePlaceholdersForRegex(templateEscapedRegex, urlAddress) ||
		containsTemplatePlaceholdersForRegex(templateUnescapedRegex, urlAddress)
}

func containsTemplatePlaceholdersForRegex(regex string, urlAddress string) bool {
	var match, _ = regexp.MatchString(regex, urlAddress)
	return match
}
