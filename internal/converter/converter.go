package converter

import (
	"io"
	"runtime"
)

// ConversionOptions holds configuration for CSV to JSON conversion
type ConversionOptions struct {
	Delimiter    rune
	HasHeader    bool
	OutputFormat string // "array" or "object"
	PrettyPrint  bool
	InferTypes   bool
}

// DefaultOptions returns default conversion options
func DefaultOptions() ConversionOptions {
	return ConversionOptions{
		Delimiter:    ',',
		HasHeader:    true,
		OutputFormat: "array",
		PrettyPrint:  true,
		InferTypes:   true,
	}
}

// ConvertCSVToJSON converts CSV data to JSON format using ultra-optimized implementation
func ConvertCSVToJSON(reader io.Reader, options ConversionOptions) ([]byte, error) {
	// Use ultra-optimized version with best performance settings
	ultraOptions := UltraOptimizedOptions{
		OptimizedConversionOptions: OptimizedConversionOptions{
			ConversionOptions: options,
			Workers:          runtime.NumCPU(),
			BatchSize:        1000,
			Streaming:        true,
		},
		UseMemoryPools: true,
		StreamingJSON:  true,
		SIMDEnabled:    true,
	}
	
	return ConvertCSVToJSONUltra(reader, ultraOptions)
}

