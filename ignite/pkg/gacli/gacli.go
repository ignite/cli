package gacli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	endpoint = "https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s"
)

type (
	// Client is an analytics client.
	Client struct {
		id     string // Google Analytics measurement ID.
		secret string // Google Analytics API secret.
	}
	// Body analytics metrics body.
	Body struct {
		ClientId string  `json:"client_id"`
		Events   []Event `json:"events"`
	}
	// Event analytics event.
	Event struct {
		Name   string `json:"name"`
		Params Metric `json:"params"`
	}
	// Metric represents a data point.
	Metric struct {
		Status             string `json:"status,omitempty"`
		OS                 string `json:"os,omitempty"`
		Arch               string `json:"arch,omitempty"`
		FullCmd            string `json:"full_command,omitempty"`
		Cmd                string `json:"command,omitempty"`
		Error              string `json:"error,omitempty"`
		Version            string `json:"version,omitempty"`
		SessionId          string `json:"session_id,omitempty"`
		EngagementTimeMsec string `json:"engagement_time_msec,omitempty"`
	}
)

// New creates a new analytics client with
// measure id and secret key.
func New(id, secret string) Client {
	return Client{
		secret: secret,
		id:     id,
	}
}

// Send sends metrics to analytics.
func (c Client) Send(body Body) error {
	// encode body
	encoded, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create an HTTP request with the payload
	url := fmt.Sprintf(endpoint, c.id, c.secret)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Error sending event. Status code: %d\n", resp.StatusCode)
	}
	return nil
}

func (c Client) SendMetric(metric Metric) error {
	metric.EngagementTimeMsec = "100"
	return c.Send(Body{
		ClientId: metric.SessionId,
		Events: []Event{{
			Name:   metric.Cmd,
			Params: metric,
		}},
	})
}
