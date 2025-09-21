#!/bin/bash

# Ultra-Optimized CSV2JSON Benchmark
# Tests the new memory pools, SIMD, and streaming optimizations

set -e

echo "⚡ Ultra-Optimized CSV2JSON Performance Test"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Build ultra-optimized version
echo -e "${BLUE}Building ultra-optimized csv2json...${NC}"
go build -o csv2json-ultra

# Create versions for comparison
echo -e "${BLUE}Creating comparison versions...${NC}"

# Legacy version (single-threaded)
cat > csv2json-legacy-only.go << 'EOF'
package main

import (
	"csv2json/internal/converter"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: csv2json-legacy-only -i input.csv -o output.json")
		os.Exit(1)
	}
	
	inputFile := os.Args[2]
	outputFile := os.Args[4]
	
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	options := converter.DefaultOptions()
	result, err := converter.ConvertCSVToJSONLegacy(file, options)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	err = os.WriteFile(outputFile, result, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Legacy conversion completed: %s -> %s\n", inputFile, outputFile)
}
EOF

go build -o csv2json-legacy-only csv2json-legacy-only.go

echo -e "${GREEN}All versions built successfully!${NC}"
echo ""

cd benchmarks

echo -e "${BLUE}Running ultra-performance comparison...${NC}"
echo ""

# Test 1: 1M rows - The ultimate performance test
echo -e "${YELLOW}🚀 Ultimate Performance Test (1M rows)${NC}"
echo "Testing all optimizations against csvtojson"
echo ""

hyperfine --warmup 1 --runs 3 \
    --export-markdown ../ultra_benchmark_results.md \
    --export-json ../ultra_benchmark_results.json \
    '../csv2json-legacy-only -i test_1m.csv -o legacy_1m.json' \
    '../csv2json-ultra -i test_1m.csv -o ultra_1m.json' \
    'csvtojson test_1m.csv > node_1m.json'

echo ""

# Memory usage comparison
echo -e "${YELLOW}📊 Memory Efficiency Test${NC}"
echo ""

echo "Memory usage comparison (1M rows):"
echo ""

echo -e "${BLUE}Legacy (single-threaded):${NC}"
/usr/bin/time -l ../csv2json-legacy-only -i test_1m.csv -o legacy_memory.json 2>&1 | grep "maximum resident set size" | awk '{printf "Memory: %.1f MB\n", $1/1024/1024}'

echo -e "${BLUE}Ultra-optimized (memory pools + SIMD + streaming):${NC}"
/usr/bin/time -l ../csv2json-ultra -i test_1m.csv -o ultra_memory.json 2>&1 | grep "maximum resident set size" | awk '{printf "Memory: %.1f MB\n", $1/1024/1024}'

echo -e "${BLUE}Node.js csvtojson:${NC}"
/usr/bin/time -l csvtojson test_1m.csv > node_memory.json 2>&1 | grep "maximum resident set size" | awk '{printf "Memory: %.1f MB\n", $1/1024/1024}'

echo ""

# CPU utilization test
echo -e "${YELLOW}🔥 CPU Utilization Test${NC}"
echo "Testing multi-core usage efficiency..."
echo ""

echo "Ultra-optimized version CPU usage:"
echo "Expected: >150% CPU (multi-core utilization)"
time ../csv2json-ultra -i test_1m.csv -o cpu_test.json

echo ""

# Throughput calculation
echo -e "${YELLOW}📈 Throughput Analysis${NC}"
echo ""

# Extract timing from hyperfine results
if [ -f "../ultra_benchmark_results.json" ]; then
    echo "Performance Summary:"
    echo "==================="
    
    # Parse JSON results (simplified)
    echo "1M rows processing times:"
    echo "- Legacy: $(grep -A 10 "csv2json-legacy-only" ../ultra_benchmark_results.md | grep "Mean" | awk '{print $4}')"
    echo "- Ultra-optimized: $(grep -A 10 "csv2json-ultra" ../ultra_benchmark_results.md | grep "Mean" | awk '{print $4}')"
    echo "- Node.js csvtojson: $(grep -A 10 "csvtojson" ../ultra_benchmark_results.md | grep "Mean" | awk '{print $4}')"
fi

echo ""

# File size verification
echo -e "${BLUE}Output Verification${NC}"
echo ""

echo "Generated file sizes:"
ls -lh *_1m.json | awk '{print $5, $9}' | while read size file; do
    echo "- $file: $size"
done

echo ""

# Correctness verification
echo "JSON validity check:"
for file in legacy_1m.json ultra_1m.json node_1m.json; do
    if [ -f "$file" ]; then
        echo -n "- $file: "
        if jq empty "$file" 2>/dev/null; then
            echo -e "${GREEN}✓ Valid${NC}"
        else
            echo -e "${RED}✗ Invalid${NC}"
        fi
    fi
done

echo ""

# Record count verification
echo "Record count verification:"
for file in legacy_1m.json ultra_1m.json node_1m.json; do
    if [ -f "$file" ]; then
        count=$(jq length "$file" 2>/dev/null || echo "Error")
        echo "- $file: $count records"
    fi
done

cd ..

echo ""
echo -e "${GREEN}🎉 Ultra-optimization benchmark completed!${NC}"
echo ""
echo -e "${PURPLE}Key Optimizations Tested:${NC}"
echo "✅ Memory pools for object reuse"
echo "✅ SIMD-style string operations"
echo "✅ Streaming JSON writer"
echo "✅ Multi-core parallel processing"
echo "✅ Optimized type parsing"
echo ""
echo -e "${YELLOW}Expected Improvements:${NC}"
echo "• 2-4x faster type parsing (SIMD optimizations)"
echo "• 30-50% lower memory usage (memory pools)"
echo "• Better scaling with CPU cores"
echo "• Reduced GC pressure"
echo ""
echo "Check ultra_benchmark_results.md for detailed analysis!"
