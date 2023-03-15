package util_test

import (
	"micro/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilSliceContains(t *testing.T) {
	sliceStr := []string{"en", "id"}

	assert.EqualValues(t, util.SliceContains(sliceStr, "en"), true)
	assert.EqualValues(t, util.SliceContains(sliceStr, "es"), false)
}

func TestUtilSentenceCase(t *testing.T) {
	t.Run("if the given string is not an empty string should return sentence cased of the given string", func(t *testing.T) {
		message := "Not Found"

		assert.Equal(t, "Not found", util.SentenceCase(message))
	})

	t.Run("if the given string is an empty string should return an empty string too", func(t *testing.T) {
		message := ""

		assert.Equal(t, "", util.SentenceCase(message))
	})
}

func TestUtilMakePathWithPrefix(t *testing.T) {
	t.Run("if the given prefix is not empty string should return path with prefix separated by back slice", func(t *testing.T) {
		prefix := "prefix"
		path := "path"

		assert.Equal(t, "prefix/path", util.MakePathWithPrefix(prefix, path))
	})

	t.Run("if the given string is an empty string should return path without prefix and back slice", func(t *testing.T) {
		prefix := ""
		path := "path"

		assert.Equal(t, "path", util.MakePathWithPrefix(prefix, path))
	})
}

func TestCensorData(t *testing.T) {
	t.Run("if given data type phone should return censored phone", func(t *testing.T) {
		clearPhone := "+6285725833220"
		censoredPhone := "+62857*****220"

		assert.Equal(t, censoredPhone, util.CensorData("phone", clearPhone))
	})

	t.Run("if given data type email should return censored email", func(t *testing.T) {
		clearEmail := "trisna.x2@gmail.com"
		censoredEmail := "tr*****x2@gmail.com"

		assert.Equal(t, censoredEmail, util.CensorData("email", clearEmail))
	})
}

func TestCensorMapStringInterfaceData(t *testing.T) {
	m := make(map[string]interface{})
	m1 := make(map[string]interface{})
	m2 := make(map[string]interface{})

	m["username"] = "rudi_tabuti"
	m1["email"] = "ruditabuti@mailinator.com"
	m1["phone"] = "+6281234567890"
	m["data"] = m1
	m2["email"] = "ruditabuti@mailinator.com"
	m2["phone"] = "+6281234567890"
	m1["data"] = m2
	m3 := util.CensorMapStringInterfaceData(m)

	assert.Equal(t, "**CENSORED**", m3["username"])
	assert.Equal(t, "ru******ti@mailinator.com", m3["data"].(map[string]interface{})["email"])
	assert.Equal(t, "+62812*****890", m3["data"].(map[string]interface{})["data"].(map[string]interface{})["phone"])
}
