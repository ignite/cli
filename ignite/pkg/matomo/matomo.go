// Package matomo is a client for Matomo to send data points for hint-type=event.
package matomo

import (
	"fmt"
	"io"
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
		idSite     uint   // Matomo ID Site.
		tokenAuth  string // Matomo Token Auth.
		httpClient http.Client
	}
	// Params analytics metrics body.
	Params struct {
		IDSite       uint   `url:"idsite"`
		Rec          uint   `url:"rec"`
		APIVersion   uint   `url:"apiv,omitempty"`
		TokenAuth    string `url:"token_auth,omitempty"`
		CustomAction uint   `url:"ca,omitempty"`
		Rand         uint64 `url:"rand,omitempty"`
		URL          string `url:"url,omitempty"`
		UTMSource    string `url:"utm_source,omitempty"`
		UTMMedium    string `url:"utm_medium,omitempty"`
		UTMCampaign  string `url:"utm_campaign,omitempty"`
		UTMContent   string `url:"utm_content,omitempty"`
		UserID       string `url:"uid,omitempty"`
		UserAgent    string `url:"ua,omitempty"`
		ActionName   string `url:"action_name"`
		Hour         int    `url:"h,omitempty"`
		Minute       int    `url:"m,omitempty"`
		Second       int    `url:"s,omitempty"`

		// Dimension1 development mode boolean.
		// 1 = devMode ON | 0 = devMode OFF.
		Dimension1 uint `url:"dimension1"`

		// Dimension2 internal boolean.
		// 1 = internal ON not supported at present | 0 = internal OFF.
		Dimension2 uint `url:"dimension2"`

		// Dimension3 is gitpod (0 or 1).
		// 1 = isGitpod ON | 0 = isGitpod OFF.
		Dimension3 uint `url:"dimension3"`

		// Dimension4 ignite version
		Dimension4 string `url:"dimension4"`

		// Dimension6 ignite config version
		Dimension6 string `url:"dimension6"`

		// Dimension7 full cli command
		Dimension7 string `url:"dimension7"`

		// Dimension11 scaffold customization type
		Dimension11 string `url:"dimension11"`

		// Dimension13 command level 1.
		Dimension13 string `url:"dimension13"`

		// Dimension14 command level 2.
		Dimension14 string `url:"dimension14"`

		// Dimension15 command level 3.
		Dimension15 string `url:"dimension15"`

		// Dimension16 command level 4.
		Dimension16 string `url:"dimension16"`

		// Dimension17 cosmos-sdk version.
		Dimension17 string `url:"dimension17"`

		// Dimension18 operational system.
		Dimension18 string `url:"dimension18"`

		// Dimension19 system architecture.
		Dimension19 string `url:"dimension19"`

		// Dimension20 golang version.
		Dimension20 string `url:"dimension20"`

		// Dimension21 command level 5.
		Dimension21 string `url:"dimension21"`

		// Dimension22 command level 6.
		Dimension22 string `url:"dimension22"`
	}
	// Metric represents a custom data.
	Metric struct {
		Name            string
		Cmd             string
		OS              string
		Arch            string
		Version         string
		CLIVersion      string
		GoVersion       string
		SDKVersion      string
		BuildDate       string
		SourceHash      string
		ConfigVersion   string
		Uname           string
		CWD             string
		ScaffoldType    string
		BuildFromSource bool
		IsGitPod        bool
		IsCI            bool
	}
)

// Option configures code generation.
type Option func(*Client)

// WithIDSite adds an id site.
func WithIDSite(idSite uint) Option {
	return func(c *Client) {
		c.idSite = idSite
	}
}

// WithTokenAuth adds a matomo token authentication.
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
	requestURL.RawQuery = queryParams.Encode()

	// Create an HTTP request with the payload.
	resp, err := c.httpClient.Get(requestURL.String())
	if err != nil {
		return errors.Wrapf(err, "error creating HTTP request: %s", requestURL.String())
	}
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(got))

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("error to add matomo analytics metric. Status code: %d", resp.StatusCode)
	}

	return nil
}

func splitCMD(cmd string) []string {
	slipCmd := strings.Split(cmd, " ")
	result := make([]string, 6)
	for i := 0; i < len(slipCmd); i++ {
		if i > 5 {
			break
		}
		result[i] = slipCmd[i]
	}
	return result
}

func (c Client) metricURL(cmd string) string {
	return fmt.Sprintf("%s/%s", c.endpoint, strings.ReplaceAll(cmd, " ", "_"))
}

// SendMetric build the metrics and send to analytics.
func (c Client) SendMetric(sessionID string, metric Metric) error {
	var (
		now       = time.Now()
		r         = rand.New(rand.NewSource(now.Unix()))
		utmMedium = "dev"
	)
	if !metric.BuildFromSource {
		utmMedium = "binary"
	}

	slipCmd := strings.Split(metric.Cmd, " ")
	result := make([]string, 5)
	for i := 0; i < len(result); i++ {
		if i <= len(result) {
			break
		}
		result[i] = slipCmd[i]
	}

	cmd := splitCMD(metric.Cmd)

	return c.Send(Params{
		IDSite:       c.idSite,
		Rec:          1,
		APIVersion:   1,
		TokenAuth:    c.tokenAuth,
		CustomAction: 1,
		Rand:         r.Uint64(),
		URL:          c.metricURL(metric.Cmd),
		UTMSource:    c.endpoint,
		UTMMedium:    utmMedium,
		UTMCampaign:  metric.SourceHash,
		UTMContent:   metric.CLIVersion,
		UserID:       sessionID,
		UserAgent:    "Go-http-client",
		ActionName:   metric.Cmd,
		Hour:         now.Hour(),
		Minute:       now.Minute(),
		Second:       now.Second(),
		Dimension1:   0,
		Dimension2:   0,
		Dimension3:   formatBool(metric.IsGitPod),
		Dimension4:   metric.Version,
		Dimension6:   metric.ConfigVersion,
		Dimension7:   metric.Cmd,
		Dimension11:  metric.ScaffoldType,
		Dimension13:  cmd[0],
		Dimension14:  cmd[1],
		Dimension15:  cmd[2],
		Dimension16:  cmd[3],
		Dimension17:  metric.SDKVersion,
		Dimension18:  metric.OS,
		Dimension19:  metric.Arch,
		Dimension20:  metric.GoVersion,
		Dimension21:  cmd[4],
		Dimension22:  cmd[5],
	})
}

// formatBool returns "1" or "0" according to the value of b.
func formatBool(b bool) uint {
	if b {
		return 1
	}
	return 0
}
