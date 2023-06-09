package encoder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"micro/pkg/encoder"
)

type JSONString struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func TestEncoderPrettyJSONWithoutIndent(t *testing.T) {
	data := JSONString{
		Code:    200,
		Data:    nil,
		Message: "OK",
	}
	json := encoder.PrettyJSONWithoutIndent(data)
	expectedJSON := "{\"code\":200,\"data\":null,\"message\":\"OK\"}\n"

	assert.Equal(t, expectedJSON, json)
}

func TestEncoderPrettyJSONWithIndent(t *testing.T) {
	data := JSONString{
		Code:    200,
		Data:    nil,
		Message: "OK",
	}
	json := encoder.PrettyJSONWithIndent(data)
	expectedJSON := "{\n\t\"code\": 200,\n\t\"data\": null,\n\t\"message\": \"OK\"\n}\n"

	assert.Equal(t, expectedJSON, json)
}
