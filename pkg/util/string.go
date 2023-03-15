package util

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// SliceContains is a function uses to check slice of strings is contain a string.
func SliceContains(slice []string, str string) bool {
	for _, a := range slice {
		if a == str {
			return true
		}
	}

	return false
}

// SentenceCase is a function uses to transform string into sentence case format.
// Input: Hello World!, Output: Hello world!
func SentenceCase(sentence string) string {
	if sentence == "" {
		return ""
	}

	tmpString := []rune(strings.ToLower(sentence))
	tmpString[0] = unicode.ToUpper(tmpString[0])

	return string(tmpString)
}

// SpaceMap is a function
func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

const (
	// LowerAlphanumeric is a set of lowercase of alphabet (a-z) and number (0-9).
	LowerAlphanumeric string = "abcdefghijklmnopqrstuvwxyz0123456789"

	// UpperAlphanumeric is a set of uppercase of alphabet (A-Z) and number (0-9).
	UpperAlphanumeric string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Alphanumeric is a set of upper & lower case of alphabet (A-Z,a-z) and number (0-9).
	Alphanumeric string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// UpperAlphabet is a set of uppercase of alphabet (A-Z).
	UpperAlphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// LowerAlphabet is a set of lowercase of alphabet (a-z).
	LowerAlphabet string = "abcdefghijklmnopqrstuvwxyz"

	// Numeric is a number (0-9)
	Numeric string = "0123456789"
)

// RandomAlphabetic is a function
func RandomAlphabetic(n int) string {
	return RandomStringWithSample(n, Alphanumeric)
}

// RandomNumber is a function
func RandomNumber(n int) string {
	return RandomStringWithSample(n, Numeric)
}

// RandomStringWithSample is a function
func RandomStringWithSample(n int, sample string) string {
	var letters = []rune(sample)

	rand.Seed(time.Now().UnixNano())

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

// ToGenericArray is a function uses to convert slice of string into slice of interface.
func ToGenericArray(slice []string) []interface{} {
	s := make([]interface{}, len(slice))
	for i, v := range slice {
		s[i] = v
	}

	return s
}

// CensorData is a function uses to censor sensitive information based on the given text type.
//
//gocyclo:ignore
func CensorData(textType string, data interface{}) interface{} {
	switch textType {
	case "email":
		if reflect.ValueOf(data).Kind() == reflect.String {
			if len(data.(string)) == 0 {
				return data
			}

			emailSplit := strings.Split(data.(string), "@")
			if len(emailSplit) > 1 {
				emailProvider := strings.Split(data.(string), "@")[1]
				emailAddress := strings.Split(data.(string), "@")[0]
				return ObfuscateText(emailAddress, 2) + "@" + emailProvider
			}

			return ObfuscateText(data.(string), 2)
		}
	case "phone":
		if reflect.ValueOf(data).Kind() == reflect.String {
			if len(data.(string)) == 0 {
				return data
			}

			phoneCodeArea := Substr(data.(string), 0, 3)
			phoneNumber := Substr(data.(string), 3, len(data.(string)))

			return phoneCodeArea + ObfuscateText(phoneNumber, 3)
		}
	case "public_key":
		return "**CENSORED**"
	case "private_key":
		return "**CENSORED**"
	case "client_id":
		return "**CENSORED**"
	case "client_secret":
		return "**CENSORED**"
	case "username":
		return "**CENSORED**"
	case "password":
		return "**CENSORED**"
	case "confirm_password":
		return "**CENSORED**"
	case "otp":
		return "**CENSORED**"
	case "token":
		return "**CENSORED**"
	case "secret":
		return "**CENSORED**"
	case "certificate", "intermediateca":
		return "**CENSORED**"
	case "Authorization":
		return "**CENSORED**"
	case "Application-Key":
		return "**CENSORED**"
	case "Access-Token":
		return "**CENSORED**"
	}

	return data
}

// ObfuscateText is a function uses to censor sensitive information contained in the given string.
// Example: email, phone, key, secret key, etc.
func ObfuscateText(text string, padding int) string {
	if text == "" {
		return text
	}

	currentPadding := padding
	censoredTextLength := len(text) - (currentPadding * 2)

	if censoredTextLength < 2 {
		currentPadding = 0
		censoredTextLength = len(text) - (currentPadding * 2)
	}

	if censoredTextLength < 5 {
		currentPadding = 1
		censoredTextLength = len(text) - (currentPadding * 2)
	}

	if censoredTextLength > 5 && censoredTextLength < 10 {
		currentPadding = 2
		censoredTextLength = len(text) - (currentPadding * 2)
	}

	if censoredTextLength > 10 {
		currentPadding = 3
		censoredTextLength = len(text) - (currentPadding * 2)
	}

	return strings.Join([]string{text[:currentPadding], strings.Repeat("*", censoredTextLength), text[censoredTextLength+currentPadding:]}, "")
}

// CensorMapStringInterfaceData is for censoring data of map[string]interface{} using CensorData.
// This function also can censor if the map value is also a map[string]interface{}.
func CensorMapStringInterfaceData(data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		if m, ok := v.(map[string]interface{}); ok {
			data[k] = CensorMapStringInterfaceData(m)
		} else {
			data[k] = CensorData(k, v)
		}
	}
	return data
}

// Substr is a function to get character from string by start and end length of char.
func Substr(s string, start, length int) string {
	asRunes := []rune(s)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
