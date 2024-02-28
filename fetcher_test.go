package main

import (
	"testing"
)

func TestFetchDataBlockString(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid URL with block id",
			url:     "https://www.notion.so/samkingston/Daily-Notes-f1ca882898194427b92d0af12d73633a",
			want:    "f1ca882898194427b92d0af12d73633a",
			wantErr: false,
		},
		{
			name:    "short url with block id",
			url:     "notion.so/samkingston/Daily-Notes-f1ca882898194427b92d0af12d73633a",
			want:    "f1ca882898194427b92d0af12d73633a",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			want:    "",
			wantErr: true,
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
