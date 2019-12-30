package protocol

import (
	"bytes"
	"net/http"
)

//Client API client
type Client struct {
	daemonURL string
}

// CheckHealth ..
func CheckHealth(daemonURL string) bool {
	resp, err := http.Head(daemonURL + V1Path + HealthPath)
	return (err == nil && resp.StatusCode == 200)
}

// NewClient ..
func NewClient(daemonURL string) (client *Client) {
	apiClient := Client{daemonURL: daemonURL}
	return &apiClient
}

// Health check health
func (client *Client) Health() bool {
	return CheckHealth(client.daemonURL)
}

// AddJob ..
func (client *Client) AddJob(data []byte) bool {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	resp, err := http.Post(client.daemonURL+V1Path+AddJobPath, "text/json", &buffer)
	return (err == nil && resp.StatusCode == 200)
}
