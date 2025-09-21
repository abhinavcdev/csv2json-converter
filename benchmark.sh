#!/bin/bash

# CSV2JSON Ultra-Performance Benchmark Script
# Comprehensive testing of the ultra-optimized converter

set -e

echo "⚡ CSV2JSON Ultra-Performance Benchmark"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Build the ultra-optimized version
echo -e "${BLUE}Building ultra-optimized csv2json...${NC}"
go build -o csv2json

# Create benchmarks directory
mkdir -p benchmarks
cd benchmarks

# Generate test data inline (no separate file needed)
echo -e "${BLUE}Generating test data...${NC}"

generate_csv() {
    local filename=$1
    local rows=$2
    
    echo "id,name,email,age,salary,active,department,join_date" > "$filename"
    
    departments=("Engineering" "Marketing" "Sales" "HR" "Finance")
    
    for ((i=1; i<=rows; i++)); do
        name="User$i"
        email="user$i@example.com"
        age=$((RANDOM % 40 + 25))
        salary=$((RANDOM % 50000 + 50000))
        active=$((RANDOM % 2))
        dept_idx=$((RANDOM % 5))
        department=${departments[$dept_idx]}
        month=$((RANDOM % 12 + 1))
        day=$((RANDOM % 28 + 1))
        join_date=$(printf "2020-%02d-%02d" $month $day)
        
        echo "$i,$name,$email,$age,$salary,$active,$department,$join_date" >> "$filename"
    done
    
    echo "Generated $filename with $rows rows"
}

# Generate test files
generate_csv "test_10k.csv" 10000
generate_csv "test_100k.csv" 100000
generate_csv "test_1m.csv" 1000000

echo -e "${GREEN}Test data generated successfully!${NC}"
echo ""

# Performance scaling test
echo -e "${YELLOW}🚀 Ultra-Performance Scaling Test${NC}"
echo "Testing performance across file sizes with all optimizations..."
echo ""

hyperfine --warmup 2 --runs 5 \
    --export-markdown ../ultra_performance_results.md \
    --export-json ../ultra_performance_results.json \
    '../csv2json -i test_10k.csv -o ultra_10k.json' \
    '../csv2json -i test_100k.csv -o ultra_100k.json' \
    '../csv2json -i test_1m.csv -o ultra_1m.json'

echo ""

# Memory efficiency test
echo -e "${YELLOW}📊 Memory Efficiency Test${NC}"
echo "Testing memory usage with ultra-optimizations..."
echo ""

echo "Memory usage comparison:"
for size in 10k 100k 1m; do
    echo -e "${BLUE}Testing $size rows:${NC}"
    /usr/bin/time -l ../csv2json -i test_${size}.csv -o memory_${size}.json 2>&1 | \
        grep "maximum resident set size" | \
        awk -v size="$size" '{printf "%s: %.1f MB\n", size, $1/1024/1024}'
done

echo ""

# Output format performance
echo -e "${YELLOW}🔄 Output Format Optimization${NC}"
echo "Comparing array vs object formats with ultra-optimizations..."
echo ""

hyperfine --warmup 1 --runs 3 \
    --export-markdown ../format_performance_results.md \
    '../csv2json -i test_100k.csv -o array_100k.json -f array' \
    '../csv2json -i test_100k.csv -o object_100k.json -f object'

echo ""

# CPU utilization test
echo -e "${YELLOW}🔥 CPU Utilization Test${NC}"
echo "Testing multi-core utilization..."
echo ""

echo "Processing 1M rows (watch CPU usage):"
time ../csv2json -i test_1m.csv -o cpu_test.json

echo ""

# Throughput calculation
echo -e "${YELLOW}📈 Throughput Analysis${NC}"
echo ""

if [ -f "../ultra_performance_results.json" ]; then
    echo "Throughput calculations:"
    echo "- 10K rows: ~$(echo "10000 / $(jq -r '.results[0].mean' ../ultra_performance_results.json)" | bc -l | xargs printf "%.0f") rows/sec"
    echo "- 100K rows: ~$(echo "100000 / $(jq -r '.results[1].mean' ../ultra_performance_results.json)" | bc -l | xargs printf "%.0f") rows/sec"
    echo "- 1M rows: ~$(echo "1000000 / $(jq -r '.results[2].mean' ../ultra_performance_results.json)" | bc -l | xargs printf "%.0f") rows/sec"
fi

echo ""

# Output validation
echo -e "${BLUE}Validating JSON outputs...${NC}"
valid_count=0
total_count=0

for file in *.json; do
    total_count=$((total_count + 1))
    if jq empty "$file" 2>/dev/null; then
        echo -e "${GREEN}✓ $file${NC}"
        valid_count=$((valid_count + 1))
    else
        echo -e "${RED}✗ $file${NC}"
    fi
done

echo ""
echo "Validation: $valid_count/$total_count files are valid JSON"

# File size analysis
echo ""
echo -e "${BLUE}Generated file sizes:${NC}"
ls -lh *.json | awk '{print $5, $9}' | while read size file; do
    echo "- $file: $size"
done
echo -e "${BLUE}Quick Performance Summary:${NC}"
echo "1M rows converted in approximately:"
grep "csv2json -i test_1m.csv -o output_1m.json" ../benchmark_results_filesize.md | tail -1 | awk '{print $4 " ± " $6}'

cd ..

echo ""
echo -e "${GREEN}🎉 Benchmark suite completed successfully!${NC}"
echo "Check the benchmark_results_*.md files for detailed results."
