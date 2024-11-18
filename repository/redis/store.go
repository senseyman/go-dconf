package redis

import (
	"context"
	"encoding/json"
)

func (c *Client) UpdateConfig(ctx context.Context, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.cli.HSetNX(ctx, dConfKey, c.appName, jsonValue).Err()
}

func (c *Client) GetConfig(ctx context.Context, obj any) error {
	res := c.cli.HGet(ctx, dConfKey, c.appName)
	if res.Err() != nil {
		return res.Err()
	}

	return res.Scan(obj)
}
