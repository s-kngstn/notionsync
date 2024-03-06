package main

import (
	"testing"
)

func TestExtractNameFromURL(t *testing.T) {
	testCases := []struct {
		name         string
		url          string
		expectedName string
		expectError  bool
	}{
		{
			name:         "Valid URL with name",
			url:          "https://www.example.so/samkingston/Tech-standup-3aa9e2a7e3d24fe1b92f7fef71e05760",
			expectedName: "tech-standup",
			expectError:  false,
		},
		{
			name:         "Valid URL with different name",
			url:          "https://www.example.so/samkingston/Networking-4d7853a8b4834ed893df710c599a1d11",
			expectedName: "networking",
			expectError:  false,
		},
		{
			name:         "Invalid URL format",
			url:          "https://www.example.so/samkingston/",
			expectedName: "",
			expectError:  true,
		},
		{
			name:         "URL without dash",
			url:          "https://www.example.so/samkingston/GoMeeting",
			expectedName: "",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, err := ExtractNameFromURL(tc.url)
			if tc.expectError && err == nil {
				t.Errorf("%s: expected error but got none", tc.name)
			} else if !tc.expectError && err != nil {
				t.Errorf("%s: did not expect error but got: %v", tc.name, err)
			} else if name != tc.expectedName {
				t.Errorf("%s: expected name %s, got %s", tc.name, tc.expectedName, name)
			}
		})
	}
}
