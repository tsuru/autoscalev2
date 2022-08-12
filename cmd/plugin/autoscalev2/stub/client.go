package stub

import (
	"context"
	"fmt"
	"time"
)

type TriggerMetadata map[string]interface{}
type Trigger struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Metadata TriggerMetadata `json:"metadata"`
}

var stubsTriggers = []Trigger{
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
}

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

func (cli *Client) ListTriggers(c context.Context, args ListTriggersArgs) ([]Trigger, error) {
	return stubsTriggers, nil
}

type UpsertTriggerArgs struct {
	Instance string
	Name     string
	Type     string
	Metadata TriggerMetadata
}

func (cli *Client) UpsertTrigger(c context.Context, args UpsertTriggerArgs) error {
	return nil
}

type DeleteTriggerArgs struct {
	Instance string
	Name     string
}

func (cli *Client) DeleteTrigger(c context.Context, args DeleteTriggerArgs) error {
	return nil
}

type GetTriggerArgs struct {
	Instance string
	Name     string
}

func (cli *Client) GetTrigger(c context.Context, args GetTriggerArgs) (Trigger, error) {
	for _, v := range stubsTriggers {
		if v.Name == args.Name {
			return v, nil
		}
	}
	return Trigger{}, fmt.Errorf("No trigger %q found in instance %q", args.Name, args.Instance)
}
