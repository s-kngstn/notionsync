package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// UserInput defines an interface for reading user input.
type UserInput interface {
	ReadString(prompt string) string
}

// RealUserInput implements UserInput for real user input scenarios.
type RealUserInput struct{}

func (RealUserInput) ReadString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
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
