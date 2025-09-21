package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run generate_test_data.go <rows> <filename>")
		os.Exit(1)
	}

	rows, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid row count: %v\n", err)
		os.Exit(1)
	}

	filename := os.Args[2]

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"id", "name", "email", "age", "salary", "department", "active", "join_date", "score", "city"}
	writer.Write(header)

	// Sample data for generation
	names := []string{"John", "Jane", "Bob", "Alice", "Charlie", "Diana", "Eve", "Frank", "Grace", "Henry"}
	departments := []string{"Engineering", "Marketing", "Sales", "HR", "Finance", "Operations", "Legal", "IT"}
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	domains := []string{"gmail.com", "yahoo.com", "hotmail.com", "company.com", "outlook.com"}

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Generating %d rows of test data...\n", rows)
	start := time.Now()

	for i := 1; i <= rows; i++ {
		name := names[rand.Intn(len(names))]
		email := fmt.Sprintf("%s%d@%s", name, rand.Intn(1000), domains[rand.Intn(len(domains))])
		age := strconv.Itoa(22 + rand.Intn(43)) // 22-65
		salary := strconv.Itoa(30000 + rand.Intn(170000)) // 30k-200k
		department := departments[rand.Intn(len(departments))]
		active := strconv.FormatBool(rand.Float32() > 0.2) // 80% active
		joinDate := fmt.Sprintf("2020-%02d-%02d", 1+rand.Intn(12), 1+rand.Intn(28))
		score := fmt.Sprintf("%.2f", rand.Float64()*100)
		city := cities[rand.Intn(len(cities))]

		row := []string{
			strconv.Itoa(i),
			name,
			email,
			age,
			salary,
			department,
			active,
			joinDate,
			score,
			city,
		}

		writer.Write(row)

		if i%100000 == 0 {
			fmt.Printf("Generated %d rows...\n", i)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Generated %d rows in %v\n", rows, elapsed)
	fmt.Printf("File saved as: %s\n", filename)

	// Print file size
	if stat, err := file.Stat(); err == nil {
		fmt.Printf("File size: %.2f MB\n", float64(stat.Size())/(1024*1024))
	}
}
