package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/siguerts1/mpc/internal/controller"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mpc [serve] --hosts hosts.yml")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "serve":
		hostsFile := flag.NewFlagSet("serve", flag.ExitOnError)
		hostsPath := hostsFile.String("hosts", "hosts.yml", "Path to hosts.yml")
		hostsFile.Parse(args)

		controller.HostsPath = *hostsPath
		http.HandleFunc("/api/instances", controller.ServeInstances)

		fmt.Println("🚀 Serving MPC API at http://localhost:9900")
		log.Fatal(http.ListenAndServe(":9900", nil))

	default:
		fmt.Println("Unknown command:", cmd)
	}
}
