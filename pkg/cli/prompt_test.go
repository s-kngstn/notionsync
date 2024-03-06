package cli

import (
	"testing"
)

func TestPrompt(t *testing.T) {
	tests := []struct {
		name      string
		prompt    string
		want      string
		responses map[string]string
	}{
		{
			name:      "prompt for name",
			prompt:    "What is your name? ",
			want:      "Alice",
			responses: map[string]string{"What is your name? ": "Alice"},
		},
		{
			name:      "prompt for age",
			prompt:    "What is your age? ",
			want:      "30",
			responses: map[string]string{"What is your age? ": "30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUI := NewMockUserInput(tt.responses)
			if got := Prompt(mockUI, tt.prompt); got != tt.want {
				t.Errorf("Prompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
