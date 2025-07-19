package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/siguerts1/mpc/internal/controller"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mpc [status|instances] --hosts hosts.yml")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	hostsFile := flag.NewFlagSet(cmd, flag.ExitOnError)
	hostsPath := hostsFile.String("hosts", "hosts.yml", "Path to hosts.yml")
	hostsFile.Parse(args)

	switch cmd {
	case "status":
		results, err := controller.CheckStatus(*hostsPath)
		if err != nil {
			log.Fatalf("❌ Failed to check hosts: %v", err)
		}
		for _, r := range results {
			if r.Error != "" {
				fmt.Printf("❌ %s - %s\n", r.Name, r.Error)
			} else {
				fmt.Printf("✅ %s - Hostname: %s - MPCD: %s\n", r.Name, r.Hostname, r.Version)
			}
		}

	case "instances":
		results, err := controller.ListInstances(*hostsPath)
		if err != nil {
			log.Fatalf("❌ Failed to list instances: %v", err)
		}

		for _, r := range results {
			fmt.Printf("🔗 %s (%s)\n", r.HostName, r.Address)
			if r.Error != "" {
				fmt.Printf("  ❌ %s\n", r.Error)
			} else {
				for _, inst := range r.Instances {
					fmt.Printf("  🖥️  %s - %s - %s\n", inst.Name, inst.State, inst.IP)
				}
			}
		}

	default:
		fmt.Println("Unknown command:", cmd)
	}
}
