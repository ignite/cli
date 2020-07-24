// Package analyticsutil is a wrapper around Segment.io's Go client to provide
// an easier interface.
package analyticsutil

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/ilgooz/analytics-go"
)

// Client is an analytics client.
type Client struct {
	loginName string
	analytics.Client
}

// New creates a new analytics client for Segment.io with Segment's
// endpoint and access key.
func New(endpoint, key string) *Client {
	// err is for validation, can be ignored.
	client, _ := analytics.NewWithConfig(key, analytics.Config{
		Endpoint: endpoint,
		Logger:   analytics.StdLogger(log.New(ioutil.Discard, "", log.LstdFlags)),
	})
	return &Client{
		Client: client,
	}
}

// Login optionally logins user to later associate user with other analytics data.
func (c *Client) Login(name, starportVersion string) error {
	c.loginName = name
	hostname, _ := os.Hostname()
	return c.Client.Enqueue(analytics.Identify{
		UserId: name,
		Traits: analytics.NewTraits().
			SetName(hostname).
			Set("hostname", hostname).
			Set("platform", runtime.GOOS).
			Set("arch", runtime.GOARCH).
			Set("starport_version", starportVersion),
	})
}

// Track adds a new analytics data to the queue.
func (c *Client) Track(track analytics.Track) error {
	track.UserId = c.loginName
	return c.Client.Enqueue(track)
}
