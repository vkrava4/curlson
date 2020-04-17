package util

import "fmt"

type GetValidatorBuilder interface {
	AddUrl(getUrl string) GetValidatorBuilder
	AddTemplate(template string) GetValidatorBuilder
	AddThreads(threads int) GetValidatorBuilder
	AddRequestCount(requestCount int) GetValidatorBuilder

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

func (b *GetValidator) Entity() ValidatorEntity {
	return b.entity
}

func (e *ValidatorEntity) Validate() *ValidationResult {
	var result = &ValidationResult{
		valid: true,
	}

	validatePositive("threads", e.threads, result)
	validatePositive("requests per thread", e.requestCount, result)

	return result
}

func validatePositive(name string, value int, result *ValidationResult) {
	if value < 1 {
		result.valid = false
		result.messages = append(result.messages, fmt.Sprintf(MsgShouldBePositive, name, value))
	}
}
