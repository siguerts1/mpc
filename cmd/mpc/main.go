package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/siguerts1/mpc/internal/controller"
)

func main() {
	hostsFile := flag.String("hosts", "hosts.yml", "Path to hosts.yml")
	flag.Parse()

	results, err := controller.CollectHostData(*hostsFile)
	if err != nil {
		log.Fatalf("❌ Failed to collect data: %v", err)
	}

	jsonOutput, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("❌ Failed to marshal JSON: %v", err)
	}

	fmt.Println(string(jsonOutput))
}
