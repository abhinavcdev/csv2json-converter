package converter

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestConvertCSVToJSON(t *testing.T) {
	tests := []struct {
		name     string
		csvData  string
		options  ConversionOptions
		expected string
	}{
		{
			name:    "Basic CSV with header",
			csvData: "name,age,active\nJohn,30,true\nJane,25,false",
			options: DefaultOptions(),
			expected: `[
  {
    "name": "John",
    "age": 30,
    "active": true
  },
  {
    "name": "Jane",
    "age": 25,
    "active": false
  }
]`,
		},
		{
			name:    "CSV without header",
			csvData: "John,30,true\nJane,25,false",
			options: ConversionOptions{
				Delimiter:    ',',
				HasHeader:    false,
				OutputFormat: "array",
				PrettyPrint:  true,
				InferTypes:   true,
			},
			expected: `[
  {
    "column_1": "John",
    "column_2": 30,
    "column_3": true
  },
  {
    "column_1": "Jane",
    "column_2": 25,
    "column_3": false
  }
]`,
		},
		{
			name:    "Object format output",
			csvData: "name,age\nJohn,30\nJane,25",
			options: ConversionOptions{
				Delimiter:    ',',
				HasHeader:    true,
				OutputFormat: "object",
				PrettyPrint:  true,
				InferTypes:   true,
			},
			expected: `{
  "name": [
    "John",
    "Jane"
  ],
  "age": [
    30,
    25
  ]
}`,
		},
		{
			name:    "Semicolon delimiter",
			csvData: "name;age;city\nJohn;30;NYC\nJane;25;LA",
			options: ConversionOptions{
				Delimiter:    ';',
				HasHeader:    true,
				OutputFormat: "array",
				PrettyPrint:  true,
				InferTypes:   true,
			},
			expected: `[
  {
    "name": "John",
    "age": 30,
    "city": "NYC"
  },
  {
    "name": "Jane",
    "age": 25,
    "city": "LA"
  }
]`,
		},
		{
			name:    "No type inference",
			csvData: "name,age,score\nJohn,30,95.5\nJane,25,87.2",
			options: ConversionOptions{
				Delimiter:    ',',
				HasHeader:    true,
				OutputFormat: "array",
				PrettyPrint:  true,
				InferTypes:   false,
			},
			expected: `[
  {
    "name": "John",
    "age": "30",
    "score": "95.5"
  },
  {
    "name": "Jane",
    "age": "25",
    "score": "87.2"
  }
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.csvData)
			result, err := ConvertCSVToJSON(reader, tt.options)
			if err != nil {
				t.Fatalf("ConvertCSVToJSON() error = %v", err)
			}

			// Parse both result and expected as JSON for comparison
			var resultJSON, expectedJSON interface{}
			if err := json.Unmarshal(result, &resultJSON); err != nil {
				t.Fatalf("Failed to parse result JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expected), &expectedJSON); err != nil {
				t.Fatalf("Failed to parse expected JSON: %v", err)
			}

			if !reflect.DeepEqual(resultJSON, expectedJSON) {
				t.Errorf("ConvertCSVToJSON() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestEmptyCSV(t *testing.T) {
	reader := strings.NewReader("")
	result, err := ConvertCSVToJSON(reader, DefaultOptions())
	if err != nil {
		t.Fatalf("ConvertCSVToJSON() error = %v", err)
	}

	expected := "[]"
	if string(result) != expected {
		t.Errorf("ConvertCSVToJSON() = %v, want %v", string(result), expected)
	}
}
