package redis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/redis/go-redis/v9"
)

func (c *Client) UpdateConfig(ctx context.Context, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.cli.HSet(ctx, dConfKey, c.appName, jsonValue).Err()
}

func (c *Client) GetConfig(ctx context.Context, obj, defaultValue any) error {
	res := c.cli.HGet(ctx, dConfKey, c.appName)
	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			v := reflect.ValueOf(obj)
			if v.Kind() == reflect.Ptr && !v.IsNil() {
				v.Elem().Set(reflect.ValueOf(defaultValue))
			}
			return nil
		}
		return res.Err()
	}

	data, err := res.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}
