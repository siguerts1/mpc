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
	IP    string `json:"ipv4"`
}

type InstancesResult struct {
	HostName  string
	Address   string
	Instances []Instance
	Error     string
}

func ListInstances(file string) ([]InstancesResult, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var list HostList
	if err := yaml.Unmarshal(content, &list); err != nil {
		return nil, err
	}

	var results []InstancesResult

	for _, h := range list.Hosts {
		url := fmt.Sprintf("http://%s:%d/instances", h.Address, h.Port)

		client := http.Client{Timeout: 4 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			results = append(results, InstancesResult{
				HostName: h.Name,
				Address:  h.Address,
				Error:    "Connection failed",
			})
			continue
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var instances []Instance
		if err := json.Unmarshal(body, &instances); err != nil {
			results = append(results, InstancesResult{
				HostName: h.Name,
				Address:  h.Address,
				Error:    "Invalid response",
			})
			continue
		}

		results = append(results, InstancesResult{
			HostName:  h.Name,
			Address:   h.Address,
			Instances: instances,
		})
	}

	return results, nil
}
