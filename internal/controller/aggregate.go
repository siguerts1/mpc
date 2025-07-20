package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type VM struct {
	Name  string `json:"name"`
	State string `json:"state"`
	IP    string `json:"ipv4"`
}

type InstanceResult struct {
	HostName  string `json:"host"`
	Address   string `json:"address"`
	Error     string `json:"error,omitempty"`
	Instances []VM   `json:"instances"`
}

type HostEntry struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

func ListInstances(hostsPath string) ([]InstanceResult, error) {
	data, err := os.ReadFile(hostsPath)
	if err != nil {
		return nil, err
	}

	var hosts []HostEntry
	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return nil, err
	}

	var results []InstanceResult
	for _, host := range hosts {
		url := fmt.Sprintf("http://%s:9901/instances", host.Address)
		cmd := exec.Command("curl", "-s", url)

		out, err := cmd.Output()
		if err != nil {
			results = append(results, InstanceResult{
				HostName: host.Name,
				Address:  host.Address,
				Error:    err.Error(),
			})
			continue
		}

		var instances []VM
		if err := json.Unmarshal(out, &instances); err != nil {
			results = append(results, InstanceResult{
				HostName: host.Name,
				Address:  host.Address,
				Error:    "Invalid JSON from agent: " + err.Error(),
			})
			continue
		}

		results = append(results, InstanceResult{
			HostName:  host.Name,
			Address:   host.Address,
			Instances: instances,
		})
	}

	return results, nil
}
