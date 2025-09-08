package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) error {
	line := string(data)
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, ":", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid header, no colons... line: %s", line)
	}

	key := parts[0]
	if strings.HasSuffix(key, " ") {
		return fmt.Errorf("invalid key format, whitespace detected between key and colon... key: %s", key)
	}

	value := strings.TrimSpace(parts[1])
	h[key] = value

	return nil
}
