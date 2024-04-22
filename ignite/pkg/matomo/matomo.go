// Package matomo is a client for Matomo to send data points for hint-type=event.
package matomo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

type (
	// Client is an analytics client.
	Client struct {
		endpoint   string
		idSite     int    // Matomo ID Site.
		tokenAuth  string // Matomo Token Auth.
		httpClient http.Client
	}
	// Params analytics metrics body.
	Params struct {
		Rec        int    `url:"rec"`
		IDSite     int    `url:"idsite"`
		TokenAuth  string `url:"token_auth"`
		ActionName string `url:"action_name"`
		UniqueID   string `url:"_id"`
		UserID     string `url:"uid"`
		Rand       int    `url:"rand,omitempty"`
		APIVersion int    `url:"apiv,omitempty"`
		ClientID   string `url:"cid,omitempty"`
		EventName  string `url:"e_n,omitempty"`
		EventValue string `url:"e_v,omitempty"`
		Hour       int    `url:"h,omitempty"`
		Minute     int    `url:"m,omitempty"`
		Second     int    `url:"s,omitempty"`
		Metric     Metric `url:"-"`
	}
	// Metric represents a custom data.
	Metric struct {
		Name     string `url:"name,omitempty"`
		Cmd      string `url:"command,omitempty"`
		OS       string `url:"os,omitempty"`
		Arch     string `url:"arch,omitempty"`
		Version  string `url:"version,omitempty"`
		IsGitPod bool   `url:"is_git_pod,omitempty"`
		IsCI     bool   `url:"is_ci,omitempty"`
	}
)

// Option configures code generation.
type Option func(*Client)

// WithIDSite adds an id site.
func WithIDSite(idSite int) Option {
	return func(c *Client) {
		c.idSite = idSite
	}
}

// WithTokenAuth adds an matomo token authentication.
func WithTokenAuth(tokenAuth string) Option {
	return func(c *Client) {
		c.tokenAuth = tokenAuth
	}
}

// New creates a new Matomo client.
func New(endpoint string, opts ...Option) Client {
	c := Client{
		endpoint: endpoint,
		httpClient: http.Client{
			Timeout: 1500 * time.Millisecond,
		},
	}
	// apply analytics options.
	for _, o := range opts {
		o(&c)
	}
	return c
}

// Send sends metric event to analytics.
func (c Client) Send(params Params) error {
	requestURL, err := url.Parse(c.endpoint)
	if err != nil {
		return err
	}

	// encode request parameters.
	queryParams, err := query.Values(params)
	if err != nil {
		return err
	}
	customMetric, err := json.Marshal(params.Metric)
	if err != nil {
		return errors.Wrapf(err, "can't marshall metric to json: %v", params.Metric)
	}
	requestURL.RawQuery = fmt.Sprintf("%s&_cvar=%s", queryParams.Encode(), customMetric)

	// Create an HTTP request with the payload.
	resp, err := c.httpClient.Get(requestURL.String())
	if err != nil {
		return errors.Wrapf(err, "error creating HTTP request: %s", requestURL.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("error to add matomo analytics metric. Status code: %d", resp.StatusCode)
	}
	return nil
}

// SendMetric build the metrics and send to analytics.
func (c Client) SendMetric(sessionID string, metric Metric) error {
	var (
		now = time.Now()
		r   = rand.New(rand.NewSource(now.Unix()))
	)
	return c.Send(Params{
		IDSite:     c.idSite,
		TokenAuth:  c.tokenAuth,
		Rec:        1,
		APIVersion: 1,
		ClientID:   sessionID,
		UniqueID:   sessionID,
		UserID:     sessionID,
		Rand:       r.Int(),
		Hour:       now.Hour(),
		Minute:     now.Minute(),
		Second:     now.Second(),
		ActionName: metric.Cmd,
		EventName:  metric.Cmd,
		EventValue: strings.ReplaceAll(metric.Cmd, " ", "_"),
		Metric:     metric,
	})
}
