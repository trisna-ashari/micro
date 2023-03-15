package util

import (
	"fmt"
	"net/url"
	"strings"

	"gitlab.privy.id/carstensz/backend/core/message-bank/exception"
)

// EncodeQueryString is a function uses to encode query URI to string.
// output: ?a=1&b=2&c=3
func EncodeQueryString(query url.Values) string {
	encodedQuery := query.Encode()
	if len(encodedQuery) > 0 {
		encodedQuery = fmt.Sprintf("?%s", encodedQuery)
	}

	return encodedQuery
}

// ValidateURI is a function uses to check the redirect URI is contain domain or not.
func ValidateURI(baseURI string, redirectURI string) (bool, error) {
	base, errParseBase := url.Parse(baseURI)
	if errParseBase != nil {
		return false, errParseBase
	}

	redirect, errParseRedirect := url.Parse(redirectURI)
	if errParseRedirect != nil {
		return false, errParseRedirect
	}

	if !strings.HasSuffix(redirect.Host, base.Host) {
		return false, exception.ErrorTextMiddlewareOauthInvalidRedirectURI
	}

	return true, nil
}

// ParseStringToQueryString is a function to parse string to QueryStrings.
func ParseStringToQueryString(query string) (url.Values, error) {
	m := make(url.Values)
	err := parseQuery(m, query)
	return m, err
}

// MergeQueryString is a function to merge QueryStrings.
func MergeQueryString(items ...url.Values) url.Values {
	var merged = make(url.Values)
	for _, item := range items {
		for key, value := range item {
			merged[key] = value
		}
	}

	return merged
}

func parseQuery(m url.Values, query string) (err error) {
	for query != "" {
		key := query

		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}

		if key == "" {
			continue
		}

		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		unescape, _ := url.QueryUnescape(value)

		m[key] = append(m[key], unescape)
	}
	return err
}
