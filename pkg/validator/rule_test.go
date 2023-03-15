package validator_test

import (
	"micro/pkg/validator"
	"regexp"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	ozzoValidation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidatorAddRule(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule()

	assert.Equal(t, rules, &validator.ValidationRules{})
}

func TestValidatorApply(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().Apply()

	assert.IsType(t, rules, []validator.ValidationRule{})
	assert.Equal(t, rules, []validator.ValidationRule(nil))
}

func TestValidatorValidationRulesByFunc(t *testing.T) {
	// TODO TestValidatorValidationRulesByFunc
}

func TestValidatorValidationRulesWhen(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().When(true, validation.AddRule().Required()).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.WhenRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesRequired(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().Required().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.By(validator.Required()))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesEmpty(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().Empty().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Empty)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesNil(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().Nil().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Nil)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesNotNil(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().NotNil().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.NotNil)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesNilOrEmpty(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().NilOrEmpty().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.NilOrNotEmpty)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesInString(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().In("foo", "bar").Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.InRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "foo/bar",
		}})
	}
}

func TestValidatorValidationRulesInInt(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().In(1, 2, 3).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.InRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "1/2/3",
		}})
	}
}

func TestValidatorValidationRulesInBool(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().In(true, false).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.InRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "true/false",
		}})
	}
}

func TestValidatorValidationRulesNotInString(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().NotIn("foo", "bar").Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.NotInRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "foo/bar",
		}})
	}
}

func TestValidatorValidationRulesNotInInt(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().NotIn(1, 2, 3).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.NotInRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "1/2/3",
		}})
	}
}

func TestValidatorValidationRulesNotInBool(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().NotIn(true, false).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.NotInRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Options",
			Value: "true/false",
		}})
	}
}

func TestValidatorValidationRulesMin(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().MinValue(10).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.ThresholdRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Length",
			Value: 10,
		}})
	}
}

func TestValidatorValidationRulesMax(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().MaxValue(255).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.ThresholdRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Length",
			Value: 255,
		}})
	}
}

func TestValidatorValidationRulesMinLength(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().MinLength(255).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.ThresholdRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Length",
			Value: 255,
		}})
	}
}

func TestValidatorValidationRulesMaxLength(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().MaxLength(255).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.ThresholdRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Length",
			Value: 255,
		}})
	}
}

func TestValidatorValidationRulesLength(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().Length(10, 255).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.LengthRule{})
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Min",
			Value: 10,
		}, {
			Key:   "Max",
			Value: 255,
		}})
	}
}

func TestValidatorValidationRulesEqualTo(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().EqualTo("field_name", "value").Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.By(validator.EqualValue("field_name")))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Target",
			Value: "attributes.field_name",
		}})
	}
}

func TestValidatorValidationRulesIsAlpha(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsAlpha().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.Alpha)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsAlphaSpace(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsAlphaSpace().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Match(regexp.MustCompile(`[a-zA-Z\\s]*$`)))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsLowerAlphaUnderscore(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsLowerAlphaUnderscore().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Match(regexp.MustCompile(`^[a-z_]*$`)))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsAlphaNumeric(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsAlphaNumeric().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.Alphanumeric)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsAlphaNumericSpace(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsAlphaNumericSpace().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Match(regexp.MustCompile(`[a-zA-Z0-9\\s]*$`)))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsAlphaNumericSpaceAndSpecialCharacter(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsAlphaNumericSpaceAndSpecialCharacter().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Match(regexp.MustCompile(`[a-zA-Z0-9_+\-\\s]*$`)))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsDate(t *testing.T) {
	dateLayout := "2006-01-28"
	validation := validator.New()
	rules := validation.AddRule().IsDate(dateLayout).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Date(dateLayout))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Layout",
			Value: dateLayout,
		}})
	}
}

func TestValidatorValidationRulesIsTime(t *testing.T) {
	timeLayout := "2006-01-28 00:01:01"
	validation := validator.New()
	rules := validation.AddRule().IsTime(timeLayout).Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Date(timeLayout))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt{{
			Key:   "Layout",
			Value: timeLayout,
		}})
	}
}

func TestValidatorValidationRulesIsDigit(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsDigit().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, ozzoValidation.Match(regexp.MustCompile(`[0-9]*$`)))
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsUUID(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsUUID().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.UUID)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsInt(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsInt().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.Int)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsFloat(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsFloat().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.Float)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsEmail(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsEmail().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.EmailFormat)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsPhone(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsPhone().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.Digit)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsURL(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsURL().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.URL)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}

func TestValidatorValidationRulesIsJSON(t *testing.T) {
	validation := validator.New()
	rules := validation.AddRule().IsJSON().Apply()

	for _, r := range rules {
		assert.IsType(t, r.Rule, is.JSON)
		assert.Equal(t, r.RuleOpt, []validator.RuleOpt(nil))
	}
}
