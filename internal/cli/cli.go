package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var (
	ErrUnexpectedInput = fmt.Errorf("unexpected input")
)

type ValidationFunc func(string) error

// BooleanQuestion asks a question that expects a boolean answer.
// The question is printed to the writer, and the answer is read from the reader.
func BooleanQuestion(writer io.Writer, reader *bufio.Reader, label string, defaultValue bool) (bool, error) {
	fmt.Fprintf(writer, "%s [Y/N]: ", label)
	line, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue, err
	}

	s := strings.ToLower(strings.TrimSpace(line))
	switch s {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s", ErrUnexpectedInput, s)
	}
}

// StringQuestion asks a question that expects a string answer.
// The question is printed to the writer, and the answer is read from the reader.
func StringQuestion(writer io.Writer, reader *bufio.Reader, label string, defaultValue string, validations ...ValidationFunc) (string, error) {
	fmt.Fprintf(writer, "%s [%s]: ", label, defaultValue)
	line, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue, err
	}

	s := strings.TrimSpace(line)
	if s == "" {
		return defaultValue, nil
	}

	for _, v := range validations {
		if err := v(s); err != nil {
			return defaultValue, err
		}
	}
	return s, nil
}
