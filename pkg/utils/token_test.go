package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const secretKey = "secret-key"

type payload_test struct {
	name       string                 // Test name
	payload    map[string]interface{} // Payload
	payloadTag string                 // Payload Tag
}

var tests = []payload_test{
	{
		name: "simple-payload",
		payload: map[string]interface{}{
			"foo": "bar",
			"1":   2,
		},
		payloadTag: "simple",
	},
	{
		name: "nested-payload",
		payload: map[string]interface{}{
			"nestedData": map[string]interface{}{
				"value": "666",
			},
		},
		payloadTag: "nested",
	},
	{
		name: "mixed-payload",
		payload: map[string]interface{}{
			"string": "hello",
			"number": 3.14,
			"bool":   true,
		},
		payloadTag: "mixed",
	},
}

func TestGenerateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		test := tests[0]
		claims := NewJwtClaims(test.payload, test.payloadTag)
		tokenStr, err := GenerateToken(claims, "secret-key", 5*time.Minute)
		require.Nil(t, err)
		require.NotEmpty(t, tokenStr)
	})
}

func TestParseToken(t *testing.T) {
	for _, test := range tests {
		testName := fmt.Sprintf("parse-%s-test", test.name)
		t.Run(testName, func(t *testing.T) {
			claims := NewJwtClaims(test.payload, test.payloadTag)
			tokenStr, err := GenerateToken(claims, secretKey, 5*time.Minute)
			require.Nil(t, err)
			require.NotEmpty(t, tokenStr)

			parsed, err := ParseToken(tokenStr, secretKey)
			require.Nil(t, err)
			require.NotNil(t, parsed)

			parsedPayload := parsed[test.payloadTag].(map[string]interface{})
			fmt.Println(test.payload, parsedPayload)
			assertPayloadEqual(t, test.payload, parsedPayload)
		})

	}

	t.Run("expired-token", func(t *testing.T) {
		test := tests[0]
		claims := NewJwtClaims(test.payload, test.payloadTag)

		tokenStr, err := GenerateToken(claims, secretKey, -5*time.Minute)
		require.Nil(t, err)
		require.NotEmpty(t, tokenStr)

		mapClaims, err := ParseToken(tokenStr, secretKey)
		require.NotNil(t, err)
		require.Nil(t, mapClaims)
		require.Equal(t, err, ErrorTokenInvalid)
	})
}

// When unmarshaling JSON data into a Go value, the JSON parser converts all numbers to float64 by default.
// To fix this, you can use type assertion to convert the float64 value to an int before comparing it.
func assertPayloadEqual(t *testing.T, expected, actual map[string]interface{}) {
	for k, v := range expected {
		actualVal, ok := actual[k]
		require.True(t, ok, "Key not found in actual payload: %s", k)

		switch typedVal := v.(type) {
		case int:
			require.Equal(t, float64(typedVal), actualVal)
		default:
			require.Equal(t, v, actualVal)
		}
	}
}
