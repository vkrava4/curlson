package util

import (
	"bufio"
	"fmt"
	"github.com/vkrava4/curlson/app"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type GetValidatorBuilder interface {
	AddUrl(getUrl string) GetValidatorBuilder
	AddTemplate(template string) GetValidatorBuilder
	AddThreads(threads int) GetValidatorBuilder
	AddRequestCount(requestCount int) GetValidatorBuilder
	AddSleep(sleep int) GetValidatorBuilder
	AddMaxDuration(sleep int) GetValidatorBuilder

	WithAppConfiguration(conf *app.Configuration) GetValidatorBuilder

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

func (b *GetValidator) WithAppConfiguration(conf *app.Configuration) GetValidatorBuilder {
	if conf.Template == nil {
		conf.Template = &app.TemplateConfiguration{}
	}

	if conf.Logs == nil {
		conf.Logs = &app.LogConfiguration{}
	}

	b.entity.conf = conf
	return b
}

func (b *GetValidator) Entity() ValidatorEntity {
	return b.entity
}

func (e *ValidatorEntity) Validate() *ValidationResult {
	var result = &ValidationResult{
		valid: true,
		conf:  e.conf,
	}

	validatePositive("Amount of threads", e.threads, result)
	validatePositive("Amount of requests per thread", e.requestCount, result)
	validatePositiveOrZero("Delay in millis property", e.sleep, result)
	validatePositiveOrZero("Maximum execution duration property", e.maxDuration, result)

	validateUrlForTemplate(e.template, e.url, result)

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

func validateEmptyTemplate(template string, urlAddress string, result *ValidationResult) bool {
	if template == "" {
		var _, errPrepareUrl = PrepareUrl(urlAddress, "")
		if errPrepareUrl != nil {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgUrlAddressInvalidWithReason, urlAddress, errPrepareUrl.Error()))
		}
		return true
	}

	return false
}

func validateUrlForTemplate(template string, urlAddress string, result *ValidationResult) {

	if !validateEmptyTemplate(template, urlAddress, result) {
		var absTemplatePath, errAbsFile = filepath.Abs(template)
		if errAbsFile != nil {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgTemplatePathInvalidWithReason, template, errAbsFile.Error()))
			return
		}

		if fileExists(absTemplatePath) {
			var templateFile, errOpenFile = os.OpenFile(template, os.O_RDONLY, filesMode)
			if errOpenFile != nil {
				result.valid = false
				result.errMessages = append(result.errMessages, fmt.Sprintf(MsgCantOpenTemplateWithReason, template, errOpenFile.Error()))
			} else {
				var templateSize, errValidateUrl = validateUrlForExistingTemplate(templateFile, urlAddress, result.warnMessages)
				_ = templateFile.Close()

				if errValidateUrl != nil {
					result.valid = false
					result.errMessages = append(result.errMessages, errValidateUrl.Error())
				}

				if templateSize > 0 && result.valid && result.conf != nil {
					result.conf.Template.Enabled = true
					result.conf.Template.Path = absTemplatePath
					result.conf.Template.Size = templateSize

				}
			}
		} else {
			result.valid = false
			result.errMessages = append(result.errMessages, fmt.Sprintf(MsgTemplateNotFound, template))
		}
	}
}

func validateUrlForExistingTemplate(templateFile *os.File, urlAddress string, warnMessages []string) (int, error) {
	var templateSize = 0
	if !ContainsTemplatePlaceholders(urlAddress) {
		warnMessages = append(warnMessages, "")
	} else {
		var reader = bufio.NewReader(templateFile)
		for {
			var line, errReadLine = reader.ReadString(filesEndLineDelimiter)

			switch {
			case errReadLine == io.EOF:
				break

			case errReadLine != nil:
				return -1, errReadLine
			}

			if len(line) == 0 {
				break
			} else {
				line = strings.TrimSuffix(line, string(filesEndLineDelimiter))
				var _, errPrepareUrl = PrepareUrl(urlAddress, line)

				if errPrepareUrl != nil {
					return -1, errPrepareUrl
				}

				templateSize++
			}
		}
	}
	return templateSize, nil
}
