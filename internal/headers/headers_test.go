package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	// Test: Valid single header
	headers := Headers{}
	data := []byte("host: localhost:42069\r\n\r\n")
	err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])

	// Test: Valid single header with extra whitespace
	headers = Headers{}
	data = []byte("host:      localhost:42069         \r\n\r\n")
	err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])

	// Test: Valid single header with capital letters in key
	headers = Headers{}
	data = []byte("Host: localhost:42069\r\n\r\n")
	err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])

	// Test: Valid header with multiple values
	headers = Headers{"host": "samplehost:9000"}
	data = []byte("Host: localhost:42069\r\n\r\n")
	err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "samplehost:9000, localhost:42069", headers["host"])

	// Test: Invalid spacing header
	headers = Headers{}
	data = []byte("       host : localhost:42069       \r\n\r\n")
	err = headers.Parse(data)
	require.Error(t, err)

	// Test: Invalid character in header key
	headers = Headers{}
	data = []byte("H©st: localhost:42069\r\n\r\n")
	err = headers.Parse(data)
	require.Error(t, err)
}
