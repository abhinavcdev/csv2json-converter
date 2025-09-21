package cmd

import (
	"csv2json/internal/converter"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile    string
	outputFile   string
	delimiter    string
	noHeader     bool
	outputFormat string
	compact      bool
	noInferTypes bool
)

var rootCmd = &cobra.Command{
	Use:   "csv2json",
	Short: "Convert CSV files to JSON format",
	Long: `A fast and flexible CSV to JSON converter with support for various output formats.
	
Examples:
  csv2json -i input.csv -o output.json
  csv2json -i data.csv --format object --delimiter ";"
  csv2json -i file.csv --no-header --compact`,
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			fmt.Println("Error: input file is required")
			cmd.Help()
			os.Exit(1)
		}

		// Parse delimiter
		var delimiterRune rune = ','
		if delimiter != "" {
			if len(delimiter) == 1 {
				delimiterRune = rune(delimiter[0])
			} else if delimiter == "\\t" {
				delimiterRune = '\t'
			} else {
				fmt.Printf("Error: invalid delimiter '%s'\n", delimiter)
				os.Exit(1)
			}
		}

		// Set up conversion options
		options := converter.ConversionOptions{
			Delimiter:    delimiterRune,
			HasHeader:    !noHeader,
			OutputFormat: outputFormat,
			PrettyPrint:  !compact,
			InferTypes:   !noInferTypes,
		}

		// Open input file
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		// Convert CSV to JSON
		jsonData, err := converter.ConvertCSVToJSON(file, options)
		if err != nil {
			fmt.Printf("Error converting CSV to JSON: %v\n", err)
			os.Exit(1)
		}

		// Output result
		if outputFile != "" {
			err = os.WriteFile(outputFile, jsonData, 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
		} else {
			fmt.Println(string(jsonData))
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input CSV file (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output JSON file (optional, prints to stdout if not specified)")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "CSV delimiter (default: comma)")
	rootCmd.Flags().BoolVar(&noHeader, "no-header", false, "CSV file has no header row")
	rootCmd.Flags().StringVar(&outputFormat, "format", "array", "Output format: 'array' or 'object'")
	rootCmd.Flags().BoolVar(&compact, "compact", false, "Compact JSON output (no pretty printing)")
	rootCmd.Flags().BoolVar(&noInferTypes, "no-infer-types", false, "Don't infer data types, keep all values as strings")

	rootCmd.MarkFlagRequired("input")
}
