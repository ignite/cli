package gacli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	// Client is an analytics client.
	Client struct {
		endpoint      string
		measurementID string // Google Analytics measurement ID.
		apiSecret     string // Google Analytics API secret.
		httpClient    http.Client
	}
	// Body analytics metrics body.
	Body struct {
		ClientID string  `json:"client_id"`
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
		SessionID          string `json:"session_id,omitempty"`
		EngagementTimeMsec string `json:"engagement_time_msec,omitempty"`
	}
)

// Option configures code generation.
type Option func(*Client)

// WithMeasurementID adds an analytics measurement ID.
func WithMeasurementID(measurementID string) Option {
	return func(c *Client) {
		c.measurementID = measurementID
	}
}

// WithAPISecret adds an analytics API secret.
func WithAPISecret(secret string) Option {
	return func(c *Client) {
		c.apiSecret = secret
	}
}

// New creates a new analytics client with
// measure id and secret key.
func New(endpoint string, opts ...Option) Client {
	c := Client{
		endpoint: endpoint,
		httpClient: http.Client{
			Timeout: 1500 * time.Millisecond,
		},
	}
	// apply user options.
	for _, o := range opts {
		o(&c)
	}
	return c
}

// Send sends metrics to analytics.
func (c Client) Send(body Body) error {
	// encode body
	encoded, err := json.Marshal(body)
	if err != nil {
		return err
	}

	requestURL, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}
	v := requestURL.Query()
	if c.measurementID != "" {
		v.Set("measurement_id", c.measurementID)
	}
	if c.apiSecret != "" {
		v.Set("api_secret", c.apiSecret)
	}
	requestURL.RawQuery = v.Encode()

	// Create an HTTP request with the payload
	resp, err := c.httpClient.Post(requestURL.String(), "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error sending event. Status code: %d", resp.StatusCode)
	}
	return nil
}

func (c Client) SendMetric(metric Metric) error {
	metric.EngagementTimeMsec = "100"
	return c.Send(Body{
		ClientID: metric.SessionID,
		Events: []Event{{
			Name:   metric.Cmd,
			Params: metric,
		}},
	})
}
