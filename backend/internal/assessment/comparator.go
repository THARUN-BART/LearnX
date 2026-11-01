package assessment

import "strings"

type Comparator struct{}

func NewComparator() *Comparator {
	return &Comparator{}
}

// Compare compares actual program output against expected output with whitespace normalization.
func (c *Comparator) Compare(actual, expected string) bool {
	normActual := c.normalize(actual)
	normExpected := c.normalize(expected)
	return normActual == normExpected
}

func (c *Comparator) normalize(s string) string {
	// 1. Normalize line endings (\r\n -> \n)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		// Trim right whitespace for each line
		cleaned = append(cleaned, strings.TrimRight(line, " \t"))
	}

	result := strings.Join(cleaned, "\n")
	// Trim trailing whitespace and newlines from entire output
	return strings.TrimSpace(result)
}
