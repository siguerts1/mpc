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

type Host struct {
    Name    string `yaml:"name"`
    Address string `yaml:"address"`
    Port    int    `yaml:"port"`
    Token   string `yaml:"token,omitempty"`
}

type HostList struct {
    Hosts []Host `yaml:"hosts"`
}

type StatusResponse struct {
    Hostname string `json:"hostname"`
    Version  string `json:"version"`
}

type HostResult struct {
    Name     string
    Hostname string
    Version  string
    Error    string
}

func CheckStatus(file string) ([]HostResult, error) {
    content, err := os.ReadFile(file)
    if err != nil {
        return nil, err
    }

    var list HostList
    if err := yaml.Unmarshal(content, &list); err != nil {
        return nil, err
    }

    var results []HostResult

    for _, h := range list.Hosts {
        url := fmt.Sprintf("http://%s:%d/status", h.Address, h.Port)

        client := http.Client{Timeout: 3 * time.Second}
        resp, err := client.Get(url)
        if err != nil {
            results = append(results, HostResult{
                Name:  h.Name,
                Error: "Connection failed",
            })
            continue
        }

        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)

        var status StatusResponse
        if err := json.Unmarshal(body, &status); err != nil {
            results = append(results, HostResult{
                Name:  h.Name,
                Error: "Invalid response",
            })
            continue
        }

        results = append(results, HostResult{
            Name:     h.Name,
            Hostname: status.Hostname,
            Version:  status.Version,
        })
    }

    return results, nil
}
