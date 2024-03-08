package cli

import (
	"fmt"
	"strings"
)

// UserInput defines an interface for reading user input.
type UserInput interface {
	ReadString(prompt string) string
}

// RealUserInput implements UserInput for real user input scenarios.
type InputReader interface {
	ReadString(delim byte) (string, error)
}

// Now, modify RealUserInput to accept an InputReader.
type RealUserInput struct {
	Reader InputReader
}

func NewRealUserInput(reader InputReader) RealUserInput {
	return RealUserInput{Reader: reader}
}

func (rui RealUserInput) ReadString(prompt string) string {
	fmt.Print(prompt)
	input, _ := rui.Reader.ReadString('\n') // Ignoring error for brevity, handle it in real code.
	return strings.TrimSpace(input)
}

// MockUserInput implements UserInput for testing scenarios.
type MockUserInput struct {
	Responses map[string]string
}

func NewMockUserInput(responses map[string]string) *MockUserInput {
	return &MockUserInput{Responses: responses}
}

func (m MockUserInput) ReadString(prompt string) string {
	return m.Responses[prompt]
}
