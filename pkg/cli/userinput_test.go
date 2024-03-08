package cli

import (
	"testing"
)

type MockInputReader struct {
	Resp string
}

func (mir MockInputReader) ReadString(delim byte) (string, error) {
	return mir.Resp, nil
}

func TestReadString(t *testing.T) {
	expected := "test input"
	mockReader := MockInputReader{Resp: expected + "\n"}
	userInput := NewRealUserInput(mockReader)

	prompt := "Please enter something: "
	result := userInput.ReadString(prompt)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestMockUserInput_ReadString(t *testing.T) {
	tests := []struct {
		name      string
		prompt    string
		want      string
		responses map[string]string
	}{
		{
			name:      "test prompt 1",
			prompt:    "Enter your name: ",
			want:      "John Doe",
			responses: map[string]string{"Enter your name: ": "John Doe"},
		},
		{
			name:      "test prompt 2",
			prompt:    "Enter your age: ",
			want:      "30",
			responses: map[string]string{"Enter your age: ": "30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInput := NewMockUserInput(tt.responses)
			if got := mockInput.ReadString(tt.prompt); got != tt.want {
				t.Errorf("MockUserInput.ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}
