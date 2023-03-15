package validator

import (
	"micro/pkg/util"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Rule is a struct represent a rule and its rule options.
type Rule struct {
	Rule    validation.Rule
	RuleOpt []RuleOpt
}

// RuleOpt is a struct represent rule option.
type RuleOpt struct {
	Key   string
	Value interface{}
}

// RuleGroup is a struct represent definition of rule.
type RuleGroup struct {
	Field string
	Data  interface{}
	Rules []Rule
}

// ValidationRule is a struct represent validation rule.
type ValidationRule struct {
	Rule    validation.Rule
	RuleKey string
	RuleOpt []RuleOpt
}

// ValidationRules is a struct to store multiple validation rule.
type ValidationRules struct {
	Rules []ValidationRule
}

// AddRule is a constructor will initialize ValidationRules.
func (v *Validator) AddRule() *ValidationRules {
	return &ValidationRules{}
}

// Apply is a function uses to apply validation.Rule were has been stored in ValidationRules.
func (vr *ValidationRules) Apply() []ValidationRule {
	return vr.Rules
}

// ByFunc is a function to set the rule that with custom closure.
func (vr *ValidationRules) ByFunc(fn validation.RuleFunc) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.By(fn),
		RuleOpt: nil,
	})

	return vr
}

// When is a function to set the rule that current field when meet predefined condition.
func (vr *ValidationRules) When(condition bool, rules *ValidationRules) *ValidationRules {
	for _, rule := range rules.Rules {
		vr.Rules = append(vr.Rules, ValidationRule{
			Rule:    validation.When(condition, rule.Rule),
			RuleOpt: rule.RuleOpt,
		})
	}

	return vr
}

// Required is a function to set the rule that current field is required.
// A value is considered empty if
//   - integer, float: zero
//   - string, array: len() == 0
//   - slice, map: nil or len() == 0
//   - interface, pointer: nil or the referenced value is empty
func (vr *ValidationRules) Required() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.By(Required()),
		RuleOpt: nil,
	})

	return vr
}

// Empty is a function to set the rule that current field value is empty.
func (vr *ValidationRules) Empty() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.Empty,
		RuleOpt: nil,
	})

	return vr
}

// Nil is a function to set the rule that current field value is nil.
func (vr *ValidationRules) Nil() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.Nil,
		RuleOpt: nil,
	})

	return vr
}

// NotNil is a function to set the rule that current field value is not nil.
func (vr *ValidationRules) NotNil() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.NotNil,
		RuleOpt: nil,
	})

	return vr
}

// NilOrEmpty is a function to set the rule that current field value is nil or empty.
func (vr *ValidationRules) NilOrEmpty() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.NilOrNotEmpty,
		RuleOpt: nil,
	})

	return vr
}

// In is a function to set the rule that current field value must be one of slices of interfaces.
func (vr *ValidationRules) In(slice ...interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.In(slice...).Error("validation.error.must_be_in"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Options",
				Value: JoinSlice(slice),
			},
		},
	})

	return vr
}

// InWithSeparator is a function to set the rule that current field value must be one of slices of interfaces with separator.
func (vr *ValidationRules) InWithSeparator(separator string, slice ...interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.In(slice...).Error("validation.error.must_be_in"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Options",
				Value: JoinSliceWithSeparator(slice, separator),
			},
		},
	})

	return vr
}

// NotIn is a function to set the rule that current field value is not on of slices.
func (vr *ValidationRules) NotIn(slice ...interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.NotIn(slice).Error("validation.error.must_be_not_in"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Options",
				Value: JoinSlice(slice),
			},
		},
	})

	return vr
}

// Between is a function to set the rule that current field value must be between min and max value.
func (vr *ValidationRules) Between(min interface{}, max interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Min(min).Error("validation.error.must_be_value_between"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Min",
				Value: min,
			},
			{
				Key:   "Max",
				Value: max,
			},
		},
	})

	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Max(max).Error("validation.error.must_be_value_between"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Min",
				Value: min,
			},
			{
				Key:   "Max",
				Value: max,
			},
		},
	})

	return vr
}

// MinValue is a function to set the rule that current field value must be no less than the length.
func (vr *ValidationRules) MinValue(length interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Min(length).Error("validation.error.must_be_no_less_than_value"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Length",
				Value: length,
			},
		},
	})

	return vr
}

// MinLength is a function to set the rule that current field value must be no less than the length.
func (vr *ValidationRules) MinLength(length interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Min(length).Error("validation.error.must_be_no_less_than_length"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Length",
				Value: length,
			},
		},
	})

	return vr
}

// MaxValue is a function to set the rule that current field value must be no more than the length.
func (vr *ValidationRules) MaxValue(length interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Max(length).Error("validation.error.must_be_no_more_than_value"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Length",
				Value: length,
			},
		},
	})

	return vr
}

// MaxLength is a function to set the rule that current field value must be no more than the length.
func (vr *ValidationRules) MaxLength(length interface{}) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Max(length).Error("validation.error.must_be_no_more_than_length"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Length",
				Value: length,
			},
		},
	})

	return vr
}

// Length is a function to set the rule that current field value must be between min and max.
func (vr *ValidationRules) Length(min int, max int) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.Length(min, max).Error("validation.error.must_be_length_between"),
		RuleKey: "length",
		RuleOpt: []RuleOpt{
			{
				Key:   "Min",
				Value: min,
			},
			{
				Key:   "Max",
				Value: max,
			},
		},
	})

	return vr
}

// EqualTo is a function to set the rule that current field value must be equal to target field's value.
func (vr *ValidationRules) EqualTo(targetField string, targetValue string) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.By(EqualValue(targetValue)),
		RuleOpt: []RuleOpt{
			{
				Key:   "Target",
				Value: "attributes." + targetField,
			},
		},
	})

	return vr
}

// IsAlpha is a function to set the rule that current field value must be letters only.
func (vr *ValidationRules) IsAlpha() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.Alpha.Error("validation.error.must_be_alpha"),
		RuleOpt: nil,
	})

	return vr
}

// IsAlphaSpace is a function to set the rule that current field value must be letters and space character only.
func (vr *ValidationRules) IsAlphaSpace() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z\s]*$`)).
			Error("validation.error.must_be_alpha_space"),
		RuleOpt: nil,
	})

	return vr
}

// IsLowerAlphaUnderscore is a function to set the rule that current field value must be lower letters and underscore character only.
func (vr *ValidationRules) IsLowerAlphaUnderscore() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-z_]*$`)).
			Error("validation.error.must_be_lower_alpha_underscore"),
		RuleOpt: nil,
	})

	return vr
}

// IsLowerAlphaUnderscoreDot is a function to set the rule that current field value must be lower letters, underscore and dot character only.
func (vr *ValidationRules) IsLowerAlphaUnderscoreDot() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-z_]*$`)).
			Error("validation.error.must_be_lower_alpha_underscore_dot"),
		RuleOpt: nil,
	})

	return vr
}

// IsNumber is a function to set the rule that current field value must be number only.
// coverage: "3","-3","0","0.0","1.0","0.1","0.0001","-555","94549870965"
func (vr *ValidationRules) IsNumber() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[+\-]?(?:(?:0|[0-9]\d*)(?:\.\d*)?|\.\d+)$`)).
			Error("validation.error.must_be_number"),
		RuleOpt: nil,
	})

	return vr
}

// IsDigit is a function to set the rule that current field value must be digit only.
func (vr *ValidationRules) IsDigit() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[0-9]*$`)).
			Error("validation.error.must_be_digit"),
		RuleOpt: nil,
	})

	return vr
}

// IsAlphaNumeric is a function to set the rule that current field value must be letters and numbers only.
func (vr *ValidationRules) IsAlphaNumeric() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.Alphanumeric.Error("validation.error.must_be_alphanumeric"),
		RuleOpt: nil,
	})

	return vr
}

// IsAlphaNumericSpace is a function to set the rule that current field value must be letters, numbers, and
// space characters only.
func (vr *ValidationRules) IsAlphaNumericSpace() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z0-9\s]*$`)).
			Error("validation.error.must_be_alphanumeric_space"),
		RuleOpt: nil,
	})

	return vr
}

// IsAlphaNumericSpaceAndSpecialCharacter is a function to set the rule that current field value must be letters, numbers,
// space characters, and allowed special characters only:
//   - dot
//   - underscore
//   - plus
func (vr *ValidationRules) IsAlphaNumericSpaceAndSpecialCharacter() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._+&#()\-\s]*$`)).
			Error("validation.error.must_be_alphanumeric_space_special_character"),
		RuleOpt: nil,
	})

	return vr
}

// IsName is a function to set the rule that current field value must an allowed name:
//   - Mr John Doe
//   - Mr. John Doe
//   - Mr J`ohn Doe
func (vr *ValidationRules) IsName() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z.,'\s]*$`)).
			Error("validation.error.must_be_valid_name"),
		RuleOpt: nil,
	})

	return vr
}

// IsValidWithCustomRegex is a function to set the rule to validate current value with custom regex.
func (vr *ValidationRules) IsValidWithCustomRegex(regex, field string) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(regex)).
			Error("validation.error.must_be_valid_" + field),
		RuleOpt: nil,
	})

	return vr
}

// IsValidWithCustomRegexRule is a function to set the rule to validate current value with custom regex.
// chars example: "a-z", "A-Z", "0-9", "_", "." "space"
func (vr *ValidationRules) IsValidWithCustomRegexRule(regex string, chars []string) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(regex)).
			Error("validation.error.must_be_contain_only"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Chars",
				Value: strings.Join(chars, " "),
			},
		},
	})

	return vr
}

// IsDepartementName is a function to set the rule that current field value must an allowed departement name:
//   - Research & Development
//   - Research & Development (R&D)
//   - Information Technology / Technology
func (vr *ValidationRules) IsDepartementName() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._+&#()\-\s]*$`)).
			Error("validation.error.must_be_valid_departement_name"),
		RuleOpt: nil,
	})

	return vr
}

// IsPositionName is a function to set the rule that current field value must an allowed position name:
//   - Chief Data Officer (CDO)
//   - Research & Development Specialist (R&D)
//   - Information Technology Support / IT Support
func (vr *ValidationRules) IsPositionName() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._+&#()\-\s]*$`)).
			Error("validation.error.must_be_valid_position_name"),
		RuleOpt: nil,
	})

	return vr
}

// IsCategoryName is a function to set the rule that current field value must an allowed category name:
//   - Software Development Project (SDO)
//   - Research & Development Specialist (R&D)
//   - Information Technology Support / IT Support
func (vr *ValidationRules) IsCategoryName() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._+&#()\-\s]*$`)).
			Error("validation.error.must_be_valid_category_name"),
		RuleOpt: nil,
	})

	return vr
}

// IsHexaColor is a function to set the rule that current field value must an allowed hexa color code:
//   - #1AFFa1
//   - #F00
//   - #F020202
func (vr *ValidationRules) IsHexaColor() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)).
			Error("validation.error.must_be_valid_hexa_color"),
		RuleOpt: nil,
	})

	return vr
}

// IsSlug is a function to set the rule that current field value is valid slug.
// Matches alphanumeric slugs without repeating dashes.
func (vr *ValidationRules) IsSlug() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^(?!-)((?:[a-z0-9]+-?)+)(?<!-)$`)).
			Error("validation.error.must_be_valid_slug"),
		RuleOpt: nil,
	})

	return vr
}

// IsFilepath is a function to set the rule that current field value is valid file path.
// Filepath format: filepath /dir/file.png or dir/file.png
func (vr *ValidationRules) IsFilepath() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^([/]?[\w.\-]+)+$`)).
			Error("validation.error.must_be_valid_filepath"),
		RuleOpt: nil,
	})

	return vr
}

// IsMimeType is a function to set the rule that current field value is valid mime type.
func (vr *ValidationRules) IsMimeType() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^([\w]+/[\w.\-]+)+$`)).
			Error("validation.error.must_be_valid_mime_type"),
		RuleOpt: nil,
	})

	return vr
}

// IsDate is a function to set the rule that current field value is valid date format.
func (vr *ValidationRules) IsDate(layout string) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Date(layout).Error("validation.error.must_be_date"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Layout",
				Value: layout,
			},
		},
	})

	return vr
}

// IsTime is a function to set the rule that current field value is valid date format.
func (vr *ValidationRules) IsTime(layout string) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Date(layout).Error("validation.error.must_be_time"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Layout",
				Value: layout,
			},
		},
	})

	return vr
}

// IsUUID is a function to set the rule that current field value must be valid UUID.
func (vr *ValidationRules) IsUUID() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.UUID.Error("validation.error.must_be_uuid"),
		RuleOpt: nil,
	})

	return vr
}

// IsInt is a function to set the rule that current field value must be int data type.
func (vr *ValidationRules) IsInt() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.Int.Error("validation.error.must_be_int"),
		RuleOpt: nil,
	})

	return vr
}

// IsFloat is a function to set the rule that current field value must be float data type.
func (vr *ValidationRules) IsFloat() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.Float.Error("validation.error.must_be_float"),
		RuleOpt: nil,
	})

	return vr
}

// IsEmail is a function to set the rule that current field value must be valid email.
func (vr *ValidationRules) IsEmail() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.EmailFormat.Error("validation.error.must_be_email"),
		RuleOpt: nil,
	})

	return vr
}

// IsPhone is a function to set the rule that current field value must be valid phone number with country code.
func (vr *ValidationRules) IsPhone() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^([+]?\d{1,2}[-\s]?|)\d{3,4}[-\s]?\d{3}[-\s]?\d{4,6}$`)).
			Error("validation.error.must_be_phone"),
		RuleOpt: nil,
	})

	return vr
}

// IsURL is a function to set the rule that current field value must be valid URL.
func (vr *ValidationRules) IsURL() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.URL.Error("validation.error.must_be_url"),
		RuleOpt: nil,
	})

	return vr
}

// IsJSON is a function to set the rule that current field value must be valid JSON.
func (vr *ValidationRules) IsJSON() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    is.JSON.Error("validation.error.must_be_json"),
		RuleOpt: nil,
	})

	return vr
}

// IsPrivyID is a function to set the rule that current field value must be valid privyID.
func (vr *ValidationRules) IsPrivyID() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`\A(DEV)?([A-Z]{2,3})\d{4}\z`)).Error("validation.error.invalid_privy_id"),
	})

	return vr
}

// IsNIK is a function to set the rule that current field value must be NIK valid value.
func (vr *ValidationRules) IsNIK() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.Match(regexp.MustCompile(`^(1[1-9]|21|[37][1-6]|5[1-3]|6[1-5]|[89][12])\d{2}\d{2}([04][1-9]|[1256][0-9]|[37][01])(0[1-9]|1[0-2])\d{2}\d{4}$`)).Error("validation.error.must_be_valid_nik"),
		RuleOpt: nil,
	})

	return vr
}

// IsPassport is a function to set the rule that current field value must be passport valid value.
func (vr *ValidationRules) IsPassport() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule:    validation.Match(regexp.MustCompile(`^[A-Za-z]{1,2}\.?\s?[0-9]\d\s?\d{3,6}[0-9A-Z]$`)).Error("validation.error.must_be_valid_passport"),
		RuleOpt: nil,
	})

	return vr
}

// IsNPWP is a function to set rule of current value must be with npwp valid value.
func (vr *ValidationRules) IsNPWP() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`^[0-9]{15,16}$`)).
			Error("validation.error.must_be_valid_npwp"),
		RuleOpt: nil,
	})

	return vr
}

// IsImagePNG is a function to set rule of current value must be with image with mimetype image/png valid value.
func (vr *ValidationRules) IsImagePNG() *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Match(regexp.MustCompile(`([^\s]+(\/(?i)(png))$)`)).
			Error("validation.error.must_be_valid_image_png"),
		RuleOpt: nil,
	})

	return vr
}

// MaxFileSize is a function to set the rule that current field value must be no more than the file size. Size value in byte.
func (vr *ValidationRules) MaxFileSize(size uint64) *ValidationRules {
	vr.Rules = append(vr.Rules, ValidationRule{
		Rule: validation.Max(size).Error("validation.error.must_be_no_more_than_size"),
		RuleOpt: []RuleOpt{
			{
				Key:   "Size",
				Value: util.ByteSize(size),
			},
		},
	})

	return vr
}
