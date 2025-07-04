package tools

import "testing"

func TestResolveAppName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "test slackbot",
			input:    "test slackbot",
			expected: "2d5a77bc-380d-49c6-bf9d-3d476fcee9dd",
		},
		{
			name:     "dev api",
			input:    "dev api",
			expected: "9ec95f78-85c7-41ce-9b7e-972ad2314bf6",
		},
		{
			name:     "prod web",
			input:    "prod web",
			expected: "prod web", // Should return original if no mapping
		},
		{
			name:     "direct app id",
			input:    "2d5a77bc-380d-49c6-bf9d-3d476fcee9dd",
			expected: "2d5a77bc-380d-49c6-bf9d-3d476fcee9dd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveAppName(tt.input)
			if result != tt.expected {
				t.Errorf("ResolveAppName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAvailableAppNames(t *testing.T) {
	names := GetAvailableAppNames()
	
	// Check that we have the expected number of mappings
	expectedCount := 16 // Based on our mapping (12 unique apps + 4 variations)
	if len(names) != expectedCount {
		t.Errorf("Expected %d app names, got %d", expectedCount, len(names))
	}
	
	// Check that some expected names are present
	expectedNames := []string{"test slackbot", "dev api", "prod slackbot"}
	for _, expectedName := range expectedNames {
		found := false
		for _, name := range names {
			if name == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected app name %q not found in available names", expectedName)
		}
	}
} 