package main

import (
	"testing"
)

func TestFetchDataBlockString(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid URL with standard UUID",
			url:     "https://www.test.com/c/e20d841f-36c8-402e-bbaf-328a2aa4247f",
			want:    "e20d841f-36c8-402e-bbaf-328a2aa4247f",
			wantErr: false,
		},
		{
			name:    "valid URL with dashed words before UUID",
			url:     "https://www.test.so/samkingston/Daily-Notes-f1ca882898194427b92d0af12d73633a",
			want:    "f1ca882898194427b92d0af12d73633a",
			wantErr: false,
		},
		{
			name:    "valid URL with compact UUID",
			url:     "https://example.com/resource/123e4567e89b12d3a456426614174000",
			want:    "123e4567e89b12d3a456426614174000",
			wantErr: false,
		},
		{
			name:    "URL without UUID",
			url:     "https://example.com/no-uuid-here",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "URL with additional path elements after UUID",
			url:     "https://www.test.so/samkingston/f1ca882898194427b92d0af12d73633a/some/other/path",
			want:    "f1ca882898194427b92d0af12d73633a",
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchDataBlockString(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDataBlockString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchDataBlockString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
