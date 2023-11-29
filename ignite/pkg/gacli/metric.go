package gacli

import (
	"strings"
)

// Metric represents a data point.
type Metric struct {
	Status  string
	FullCmd string
	Cmd     string
	Error   string
	User    string
	Version string
}

func (c Client) SendMetric(metric Metric) error {
	return c.Send(Body{
		ClientId: metric.User,
		Events: []Event{{
			Name: metric.Cmd,
			Params: Params{
				CampaignId: metric.Cmd,
				Campaign:   metric.FullCmd,
				Source:     metric.Version,
				Medium:     metric.Status,
				Term:       strings.ReplaceAll(metric.FullCmd, " ", "+"),
				Content:    metric.Error,
			},
		}},
	})
}
