package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {                               //func for testing headers parsing
	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"]) // Checked lowercase
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("Test: Case Sensitivity", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("USER-AGENT: curl/7.68.0\r\n\r\n")
		_, _, err := headers.Parse(data)
		require.NoError(t, err)
		// Should be accessible via lowercase key
		assert.Equal(t, "curl/7.68.0", headers["user-agent"])
		// Should NOT exist in original casing
		_, exists := headers["USER-AGENT"]
		assert.False(t, exists)
	})

	t.Run("Test: Invalid Character", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost\r\n\r\n")
		_, _, err := headers.Parse(data)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("Valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.True(t, done)
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host : localhost\r\n\r\n")
		_, _, err := headers.Parse(data)
		require.Error(t, err)
	})

	t.Run("Valid 2 headers with same key", func(t *testing.T) {
		headers := NewHeaders()
		// Manually add the first occurrence (normalized to lowercase)
		headers["set-person"] = "lane-loves-go"

		// Data for the second occurrence
		data := []byte("Set-Person: prime-loves-zig\r\n\r\n")
		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		// The map should now contain both values separated by a comma
		assert.Equal(t, "lane-loves-go, prime-loves-zig", headers["set-person"])
		assert.Equal(t, 29, n)
		assert.False(t, done)
	})
}
