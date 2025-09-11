package headers

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

func (h Headers) Parse(line string) error {
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, ":", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid header, no colons... line: %s", line)
	}

	key := strings.ToLower(parts[0])
	if !keyIsValid(key) {
		return fmt.Errorf("invalid key format: %s", key)
	}

	value := strings.TrimSpace(parts[1])
	current_value, exists := h[key]
	if exists {
		value = fmt.Sprintf("%s, %s", current_value, value)
	}

	h[key] = value

	return nil
}

func (h Headers) Get(key string) string {
	raw, exists := h[key]
	if !exists {
		return ""
	}

	values := strings.Split(raw, ",")
	return values[0]
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func keyIsValid(key string) bool {
	if strings.HasSuffix(key, " ") {
		return false
	}

	special_characters := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	for _, c := range key {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !slices.Contains(special_characters, c) {
			return false
		}
	}

	return true
}
