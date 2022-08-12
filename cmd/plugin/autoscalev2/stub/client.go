package stub

import (
	"context"
	"time"
)

type Client struct{}

type ClientOptions struct {
	Timeout time.Duration
}

func NewClientWithOptions(address string, user string, password string, opts ClientOptions) (Client, error) {
	return Client{}, nil
}

func NewClientThroughTsuruWithOptions(target string, token string, service string, opts ClientOptions) (Client, error) {
	return Client{}, nil
}

type ListTriggersArgs struct {
	Instance string
}

type TriggerMetadata map[string]interface{}

type UpsertTriggerArgs struct {
	Instance string
	Name     string
	Type     string
	Metadata TriggerMetadata
}

type Trigger struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Metadata TriggerMetadata `json:"metadata"`
}

func (cli *Client) ListTriggers(c context.Context, args ListTriggersArgs) ([]Trigger, error) {
	return []Trigger{
		{
			"my-trigger-cron",
			"Cron",
			TriggerMetadata{
				"timezone":        "America/Sao_Paulo",
				"start":           "30 * * * *",
				"end":             "45 * * * *",
				"desiredReplicas": "10",
			},
		},
		{
			"my-trigger-prom",
			"Prometheus",
			TriggerMetadata{
				"serverAddress":       "http://prometheus-host:9090",
				"metricName":          "http_requests_total",
				"query":               "sum(rate(http_requests_total{deployment=\"my-deployment\"}[2m]))",
				"threshold":           "100.50",
				"activationThreshold": "5.5",
			},
		},
	}, nil
}
func (cli *Client) UpsertTrigger(c context.Context, args UpsertTriggerArgs) error {
	return nil
}
