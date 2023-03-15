package util_test

import (
	"micro/pkg/util"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilEncodedQueryString(t *testing.T) {
	queryString := make(url.Values)
	queryString.Set("grant_type", "code")
	queryString.Set("redirect_uri", "http://localhost")
	encodedQueryString := util.EncodeQueryString(queryString)

	assert.EqualValues(t, "?grant_type=code&redirect_uri=http%3A%2F%2Flocalhost", encodedQueryString)
}

func TestUtilValidateURI(t *testing.T) {
	baseURI := "http://example.com"
	redirectURI := "http://oauth.example.com"
	valid, err := util.ValidateURI(baseURI, redirectURI)

	assert.NoError(t, nil, err)
	assert.True(t, valid)
}

func TestUtilParseStringToQueryString(t *testing.T) {
	strString := "equal[name]=Admin&equal[name]=User&not[name]=Guest&like[name]=Super %26 Admin&&"
	expectedQueryString := url.Values{
		"equal[name]": []string{"Admin", "User"},
		"like[name]":  []string{"Super & Admin"},
		"not[name]":   []string{"Guest"},
	}
	queryString, _ := util.ParseStringToQueryString(strString)

	assert.EqualValues(t, expectedQueryString, queryString)
}

func TestUtilMergeQueryString(t *testing.T) {
	strStringOne := "equal[name]=Admin"
	strStringTwo := "not[name]=Guest"
	strStringThree := "like[name]=Super Admin"

	expectedQueryString := url.Values{
		"equal[name]": []string{"Admin"},
		"like[name]":  []string{"Super Admin"},
		"not[name]":   []string{"Guest"},
	}
	queryStringOne, _ := util.ParseStringToQueryString(strStringOne)
	queryStringTwo, _ := util.ParseStringToQueryString(strStringTwo)
	queryStringThree, _ := util.ParseStringToQueryString(strStringThree)
	queryStringMerged := util.MergeQueryString(queryStringOne, queryStringTwo, queryStringThree)

	assert.EqualValues(t, expectedQueryString, queryStringMerged)
}
