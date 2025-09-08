package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	// Test: Valid single header
	headers := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])

	// Test: Valid single header with extra whitespace
	headers = Headers{}
	data = []byte("Host:      localhost:42069         \r\n\r\n")
	err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])

	// Test: Invalid spacing header
	headers = Headers{}
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	err = headers.Parse(data)
	require.Error(t, err)
}
