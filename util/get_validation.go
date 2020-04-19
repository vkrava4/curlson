package util

import (
	"fmt"
	"os"
	"path/filepath"
)

type GetValidatorBuilder interface {
	AddUrl(getUrl string) GetValidatorBuilder
	AddTemplate(template string) GetValidatorBuilder
	AddThreads(threads int) GetValidatorBuilder
	AddRequestCount(requestCount int) GetValidatorBuilder
	AddSleep(sleep int) GetValidatorBuilder
	AddMaxDuration(sleep int) GetValidatorBuilder

	WithAppConfiguration(conf *AppConfiguration) GetValidatorBuilder

	Entity() ValidatorEntity
}

type GetValidator struct {
	entity ValidatorEntity
}

func (b *GetValidator) AddUrl(getUrl string) GetValidatorBuilder {
	b.entity.url = getUrl
	return b
}

func (b *GetValidator) AddTemplate(template string) GetValidatorBuilder {
	b.entity.template = template
	return b
}

func (b *GetValidator) AddThreads(threads int) GetValidatorBuilder {
	b.entity.threads = threads
	return b
}

func (b *GetValidator) AddRequestCount(requestCount int) GetValidatorBuilder {
	b.entity.requestCount = requestCount
	return b
}

func (b *GetValidator) AddSleep(sleep int) GetValidatorBuilder {
	b.entity.sleep = sleep
	return b
}

func (b *GetValidator) AddMaxDuration(maxDuration int) GetValidatorBuilder {
	b.entity.maxDuration = maxDuration
	return b
}

func (b *GetValidator) WithAppConfiguration(conf *AppConfiguration) GetValidatorBuilder {
	b.entity.conf = conf
	return b
}

func (b *GetValidator) Entity() ValidatorEntity {
	return b.entity
}

func (e *ValidatorEntity) Validate() *ValidationResult {
	var result = &ValidationResult{
		valid: true,
	}

	validatePositive("Amount of threads", e.threads, result)
	validatePositive("Amount of requests per thread", e.requestCount, result)
	validatePositiveOrZero("Delay in millis property", e.sleep, result)
	validatePositiveOrZero("Maximum execution duration property", e.maxDuration, result)

	validateTemplate(e.template, e.url, result, e.conf)

	return result
}

func validatePositive(description string, value int, result *ValidationResult) {
	if value < 1 {
		result.valid = false
		result.errMessages = append(result.errMessages, fmt.Sprintf(MsgShouldBePositive, description, value))
	}
}

func validatePositiveOrZero(description string, value int, result *ValidationResult) {
	if value < 0 {
		result.valid = false
		result.errMessages = append(result.errMessages, fmt.Sprintf(MsgShouldBePositiveOrZero, description, value))
	}
}

func validateUrl(urlAddress string, result *ValidationResult) {

}

func validateTemplate(template string, url string, result *ValidationResult, conf *AppConfiguration) {
	if template == "" {
		if ContainsTemplatePlaceholders(url) {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgUrlAddressInvalidWithReason, url, "URL address contains placeholder(s) for missing template file"))
			return
		}
	} else {
		var absTemplatePath, errAbsFile = filepath.Abs(template)
		if errAbsFile != nil {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgTemplatePathInvalidWithReason, template, errAbsFile.Error()))
			return
		}

		if FileExists(absTemplatePath) {
			var templateFile, errOpenFile = os.OpenFile(template, os.O_RDONLY, defaultMode)
			if errOpenFile != nil {
				result.valid = false
				result.errMessages = append(result.errMessages, fmt.Sprintf(MsgCantOpenTemplateWithReason, template, errOpenFile.Error()))
				return
			} else {
				// TODO validate whether ALL file lines match to URL

				_ = templateFile.Close()
				if !ContainsTemplatePlaceholders(url) {
					result.warnMessages = append(result.warnMessages, "")
				}

				if conf != nil {
					conf.templatingEnabled = true
				}
			}
		} else {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgTemplateNotFound, template))
			return
		}
	}
}
