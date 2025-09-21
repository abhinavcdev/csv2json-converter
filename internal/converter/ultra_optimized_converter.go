package converter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
)

// Memory pools for object reuse
var (
	objectPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 16)
		},
	}
	
	slicePool = sync.Pool{
		New: func() interface{} {
			return make([]interface{}, 0, 32)
		},
	}
)

// OptimizedConversionOptions holds configuration for optimized CSV to JSON conversion
type OptimizedConversionOptions struct {
	ConversionOptions
	Workers   int
	BatchSize int
	Streaming bool
}

// UltraOptimizedOptions extends OptimizedConversionOptions with additional ultra-optimizations
type UltraOptimizedOptions struct {
	OptimizedConversionOptions
	UseMemoryPools bool
	StreamingJSON  bool
	SIMDEnabled    bool
}

// DefaultOptimizedOptions returns default optimized conversion options
func DefaultOptimizedOptions() OptimizedConversionOptions {
	return OptimizedConversionOptions{
		ConversionOptions: DefaultOptions(),
		Workers:          runtime.NumCPU(),
		BatchSize:        1000,
		Streaming:        true,
	}
}

// DefaultUltraOptimizedOptions returns ultra-optimized settings
func DefaultUltraOptimizedOptions() UltraOptimizedOptions {
	return UltraOptimizedOptions{
		OptimizedConversionOptions: DefaultOptimizedOptions(),
		UseMemoryPools:            true,
		StreamingJSON:             true,
		SIMDEnabled:              true,
	}
}

// ConvertCSVToJSONUltra converts CSV data using ultra-optimized implementation
func ConvertCSVToJSONUltra(reader io.Reader, options UltraOptimizedOptions) ([]byte, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = options.ConversionOptions.Delimiter

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) == 0 {
		return []byte("[]"), nil
	}

	var headers []string
	var dataRows [][]string

	if options.ConversionOptions.HasHeader {
		headers = records[0]
		dataRows = records[1:]
	} else {
		// Generate generic headers
		for i := 0; i < len(records[0]); i++ {
			headers = append(headers, fmt.Sprintf("column_%d", i+1))
		}
		dataRows = records
	}

	if options.ConversionOptions.OutputFormat == "object" {
		return convertToObjectUltra(dataRows, headers, options)
	}
	
	return convertToArrayUltra(dataRows, headers, options)
}

// convertToArrayUltra uses ultra-optimizations for array format
func convertToArrayUltra(dataRows [][]string, headers []string, options UltraOptimizedOptions) ([]byte, error) {
	// Use worker pools for parallel processing
	numWorkers := options.Workers
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	rowChan := make(chan int, len(dataRows))
	resultChan := make(chan map[string]interface{}, len(dataRows))
	
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rowIdx := range rowChan {
				obj := processRowUltra(dataRows[rowIdx], headers, options)
				resultChan <- obj
			}
		}()
	}

	// Send work
	go func() {
		defer close(rowChan)
		for i := range dataRows {
			rowChan <- i
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	jsonArray := make([]map[string]interface{}, 0, len(dataRows))
	for obj := range resultChan {
		jsonArray = append(jsonArray, obj)
	}

	// Use fast JSON marshaling
	if options.ConversionOptions.PrettyPrint {
		return json.MarshalIndent(jsonArray, "", "  ")
	}
	return json.Marshal(jsonArray)
}

// convertToObjectUltra uses ultra-optimizations for object format
func convertToObjectUltra(dataRows [][]string, headers []string, options UltraOptimizedOptions) ([]byte, error) {
	jsonObj := make(map[string]interface{})
	
	// Process columns in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for i, header := range headers {
		wg.Add(1)
		go func(colIdx int, colName string) {
			defer wg.Done()
			
			var column []interface{}
			if options.UseMemoryPools {
				slice := slicePool.Get().([]interface{})
				column = slice[:0] // Reset length but keep capacity
				defer slicePool.Put(column)
			} else {
				column = make([]interface{}, 0, len(dataRows))
			}
			
			for _, row := range dataRows {
				if colIdx < len(row) {
					value := parseValueUltra(row[colIdx], options.ConversionOptions.InferTypes, options.SIMDEnabled)
					column = append(column, value)
				} else {
					column = append(column, nil)
				}
			}
			
			// Copy column data before returning to pool
			finalColumn := make([]interface{}, len(column))
			copy(finalColumn, column)
			
			mu.Lock()
			jsonObj[colName] = finalColumn
			mu.Unlock()
		}(i, header)
	}
	
	wg.Wait()

	if options.ConversionOptions.PrettyPrint {
		return json.MarshalIndent(jsonObj, "", "  ")
	}
	return json.Marshal(jsonObj)
}

// processRowUltra processes a single row with ultra-optimizations
func processRowUltra(row []string, headers []string, options UltraOptimizedOptions) map[string]interface{} {
	var obj map[string]interface{}
	
	if options.UseMemoryPools {
		obj = objectPool.Get().(map[string]interface{})
		// Clear the map but keep capacity
		for k := range obj {
			delete(obj, k)
		}
		defer objectPool.Put(obj)
	} else {
		obj = make(map[string]interface{}, len(headers))
	}
	
	// Create result map (copy from pooled object)
	result := make(map[string]interface{}, len(headers))
	
	for j, header := range headers {
		if j < len(row) {
			value := parseValueUltra(row[j], options.ConversionOptions.InferTypes, options.SIMDEnabled)
			result[header] = value
		} else {
			result[header] = nil
		}
	}
	
	return result
}

// parseValueUltra provides ultra-fast type inference with SIMD-style optimizations
func parseValueUltra(s string, inferTypes bool, simdEnabled bool) interface{} {
	if !inferTypes {
		return s
	}
	
	if len(s) == 0 {
		return nil
	}
	
	if simdEnabled {
		return parseValueSIMD(s)
	}
	
	return parseValueFast(s)
}

// parseValueSIMD uses SIMD-style optimizations for type parsing
func parseValueSIMD(s string) interface{} {
	// Fast path for common cases
	switch len(s) {
	case 0:
		return nil
	case 1:
		if s[0] >= '0' && s[0] <= '9' {
			return int64(s[0] - '0')
		}
	case 4:
		if s == "true" {
			return true
		}
	case 5:
		if s == "false" {
			return false
		}
	}
	
	// Fast number detection using first character
	if len(s) > 0 && (s[0] >= '0' && s[0] <= '9' || s[0] == '-' || s[0] == '+') {
		// Try integer first
		if intVal, err := strconv.ParseInt(s, 10, 64); err == nil {
			return intVal
		}
		// Try float
		if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
			return floatVal
		}
	}
	
	return s
}

// parseValueFast provides fast type inference
func parseValueFast(s string) interface{} {
	if len(s) == 0 {
		return nil
	}
	
	// Try integer
	if intVal, err := strconv.ParseInt(s, 10, 64); err == nil {
		return intVal
	}
	
	// Try float
	if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		return floatVal
	}
	
	// Try boolean
	if boolVal, err := strconv.ParseBool(s); err == nil {
		return boolVal
	}
	
	return s
}
