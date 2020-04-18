package util

import "fmt"

type GetValidatorBuilder interface {
	AddUrl(getUrl string) GetValidatorBuilder
	AddTemplate(template string) GetValidatorBuilder
	AddThreads(threads int) GetValidatorBuilder
	AddRequestCount(requestCount int) GetValidatorBuilder
	AddSleep(sleep int) GetValidatorBuilder
	AddMaxDuration(sleep int) GetValidatorBuilder

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

	return result
}

func validatePositive(description string, value int, result *ValidationResult) {
	if value < 1 {
		result.valid = false
		result.messages = append(result.messages, fmt.Sprintf(MsgShouldBePositive, description, value))
	}
}

func validatePositiveOrZero(description string, value int, result *ValidationResult) {
	if value < 0 {
		result.valid = false
		result.messages = append(result.messages, fmt.Sprintf(MsgShouldBePositiveOrZero, description, value))
	}
}
