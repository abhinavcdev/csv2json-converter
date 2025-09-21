package main

import (
	"csv2json/cmd"
	"csv2json/internal/api"
	"fmt"
	"os"
)

func main() {
	// Check if running as API server
	if len(os.Args) > 1 && os.Args[1] == "-server" {
		// Start API server
		fmt.Println("Starting CSV2JSON API server...")
		api.StartServer()
	} else {
		// Run CLI
		cmd.Execute()
	}
}
