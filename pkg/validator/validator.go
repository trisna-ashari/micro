package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"micro/pkg/exception"
)

// Validator is a struct represent group of rules.
type Validator struct {
	Rules []RuleGroup
	Scope string
	Mode  string
}

// New is a constructor will initialize Validator.
func New(options ...Option) *Validator {
	validator := &Validator{}

	for _, opt := range options {
		opt(validator)
	}

	return validator
}

// Set is a function uses to set a Rule to the field.
func (v *Validator) Set(field string, data interface{}, rules []ValidationRule) *Validator {
	var fieldRule []Rule
	for _, rule := range rules {
		fieldRule = append(fieldRule, Rule{
			Rule:    rule.Rule,
			RuleOpt: rule.RuleOpt,
		})
	}

	rule := RuleGroup{
		Field: field,
		Data:  data,
		Rules: fieldRule,
	}
	v.Rules = append(v.Rules, rule)

	return v
}

// Validate is a function uses to validate the value of the field
// by predefined Validator.Rules.
func (v *Validator) Validate() exception.ErrorValidators {
	var errMsg exception.ErrorValidators

	for _, ruleGroup := range v.Rules {
		var rules []validation.Rule
		var rulesOpt []RuleOpt

		for _, rule := range ruleGroup.Rules {
			rules = append(rules, rule.Rule)
			rulesOpt = append(rulesOpt, rule.RuleOpt...)
		}
		err := validation.Validate(ruleGroup.Data,
			rules...,
		)

		if err == nil {
			continue

		}

		errMsgData := make(map[string]interface{})
		errMsgData[ruleGroup.Field] = ruleGroup.Data

		for _, opt := range rulesOpt {
			errMsgData[opt.Key] = opt.Value
		}

		errMsg = append(errMsg, exception.ErrorValidator{
			Scope: v.Scope,
			Field: ruleGroup.Field,
			Msg:   err.Error(),
			Data:  errMsgData,
		})
	}

	return errMsg
}
