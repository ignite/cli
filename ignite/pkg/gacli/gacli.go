// Package gacli is a client for Google Analytics to send data points for hint-type=event.
package gacli

import (
	ga "github.com/ozgur-yalcin/google-analytics/src"
)

// Client is an analytics client.
type Client struct {
	id string // Google Analytics ID
}

// New creates a new analytics client for Segment.io with Segment's
// endpoint and access key.
func New(id string) *Client {
	return &Client{
		id: id,
	}
}

// Metric represents a data point.
type Metric struct {
	Category string
	Action   string
	Label    string
	Value    string
	User     string
	Version  string
}

// Send sends metrics to GA.
func (c *Client) Send(metric Metric) string {
	api := new(ga.API)
	api.ContentType = "application/x-www-form-urlencoded"

	client := new(ga.Client)
	client.ProtocolVersion = "1"
	client.ClientID = metric.User
	client.TrackingID = c.id
	client.HitType = "event"
	client.DocumentLocationURL = "https://github.com/ignite/cli"
	client.DocumentTitle = metric.Action
	client.DocumentEncoding = "UTF-8"
	client.EventCategory = metric.Category
	client.EventAction = metric.Action
	client.EventLabel = metric.Label
	client.EventValue = metric.Value

	return api.Send(client)
}
