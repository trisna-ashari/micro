package configurator

import (
	"os"
	"strconv"
	"strings"
)

// GetEnv is a function uses to read an environment or return a default value.
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

// GetEnvAsInt is a function uses to read an environment variable into integer or return a default value.
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// GetEnvAsBool is a function uses to read an environment variable into a bool or return default value.
func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// GetEnvAsSliceOfString is a function to read an environment variable into a string slice or return default value.
func GetEnvAsSliceOfString(name string, defaultVal []string, sep string) []string {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

// GetEnvAsSliceOfInt is a function to read an environment variable into a int slice or return default value.
func GetEnvAsSliceOfInt(name string, defaultVal []int, sep string) []int {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	var valInt []int

	val := strings.Split(valStr, sep)
	for _, v := range val {
		if value, err := strconv.Atoi(v); err == nil {
			valInt = append(valInt, value)
		}
	}

	return valInt
}

// GetEnvAsSliceOfBool is a function to read an environment variable into a bool slice or return default value.
func GetEnvAsSliceOfBool(name string, defaultVal []bool, sep string) []bool {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	var valInt []bool

	val := strings.Split(valStr, sep)
	for _, v := range val {
		if value, err := strconv.ParseBool(v); err == nil {
			valInt = append(valInt, value)
		}
	}

	return valInt
}
