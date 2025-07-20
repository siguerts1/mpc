package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Instance struct {
	Name  string `json:"name"`
	State string `json:"state"`
	IPv4  string `json:"ipv4"`
}

type HostData struct {
	HostName  string     `json:"host"`
	Address   string     `json:"address"`
	Reachable bool       `json:"reachable"`
	Hostname  string     `json:"hostname,omitempty"`
	Version   string     `json:"version,omitempty"`
	Instances []Instance `json:"instances,omitempty"`
	Error     string     `json:"error,omitempty"`
}

type hostFileEntry struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type HostsList struct {
	Hosts []hostFileEntry `yaml:"hosts"`
}

type StatusResponse struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
}

func CollectHostData(path string) ([]HostData, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var parsed HostsList
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		return nil, err
	}

	var results []HostData

	for _, host := range parsed.Hosts {
		base := fmt.Sprintf("http://%s:%d", host.Address, host.Port)
		client := http.Client{Timeout: 3 * time.Second}

		// Check /status
		statusURL := base + "/status"
		resp, err := client.Get(statusURL)
		if err != nil {
			results = append(results, HostData{
				HostName:  host.Name,
				Address:   host.Address,
				Reachable: false,
				Error:     "connection failed",
			})
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var status StatusResponse
		if err := json.Unmarshal(body, &status); err != nil {
			results = append(results, HostData{
				HostName:  host.Name,
				Address:   host.Address,
				Reachable: false,
				Error:     "invalid status response",
			})
			continue
		}

		// Fetch /instances
		instancesURL := base + "/instances"
		instResp, err := client.Get(instancesURL)
		if err != nil {
			results = append(results, HostData{
				HostName:  host.Name,
				Address:   host.Address,
				Reachable: true,
				Hostname:  status.Hostname,
				Version:   status.Version,
				Error:     "failed to get instances",
			})
			continue
		}
		defer instResp.Body.Close()

		instBody, _ := io.ReadAll(instResp.Body)
		var instances []Instance
		if err := json.Unmarshal(instBody, &instances); err != nil {
			results = append(results, HostData{
				HostName:  host.Name,
				Address:   host.Address,
				Reachable: true,
				Hostname:  status.Hostname,
				Version:   status.Version,
				Error:     "invalid instance format",
			})
			continue
		}

		results = append(results, HostData{
			HostName:  host.Name,
			Address:   host.Address,
			Reachable: true,
			Hostname:  status.Hostname,
			Version:   status.Version,
			Instances: instances,
		})
	}

	return results, nil
}
