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
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj must be a pointer")
	}

	res := c.cli.HGet(ctx, dConfKey, c.appName)
	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			objValue := reflect.ValueOf(obj).Elem()
			defaultValueValue := reflect.ValueOf(defaultValue)
			if objValue.Kind() == defaultValueValue.Kind() {
				objValue.Set(defaultValueValue)
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
