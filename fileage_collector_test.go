package main

import (
	"reflect"
	"testing"
)

type FilePatternTest struct {
	name           string
	pattern        []string
	expected       []string
	expected_error error
}

func TestGetFilePaths(t *testing.T) {
	var emptySlice []string

	tests := []FilePatternTest{
		{
			name:           "none",
			pattern:        []string{"foobar"},
			expected:       emptySlice,
			expected_error: nil,
		},
		{
			name:           "simple-pattern",
			pattern:        []string{"fileage_collector*.go"},
			expected:       []string{"fileage_collector.go", "fileage_collector_test.go"},
			expected_error: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := getFilePaths(test.pattern)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatal("actual:", actual, "expected:", test.expected)
			}
			if err != test.expected_error {
				t.Fatal("actual error:", err, "expected error:", test.expected_error)
			}
		})
	}
}

func TestGetFileStats(t *testing.T) {
	actual, err := getFileStats("./fileage_collector_test.go")
	if err != nil {
		t.Fatal("actual err:", err)
	}
	if actual.age_in_seconds < 1.0 {
		t.Fatal("actual:", actual)
	}
}
