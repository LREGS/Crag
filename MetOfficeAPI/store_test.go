package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	rdbCmd, err := startRedis()
	if err != nil {
		fmt.Printf("failed starting redis %s", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:4420"})

	code := m.Run()

	rdb.Close()
	if err := stopRedis(rdbCmd); err != nil {
		fmt.Printf("failed ending redis %s", err)
	}

	os.Exit(code)

}

func startRedis() (*exec.Cmd, error) {
	cmd := exec.Command("redis-server", "--port", "4420")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func stopRedis(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("stopping redis failed %s", err)
	}

	return nil
}

func NewClient(t *testing.T) *redis.Client {
	cmd, err := startRedis()
	if err != nil {
		t.Fatalf("failed getting rdb %s", err)
	}
	defer func() {
		if err := stopRedis(cmd); err != nil {
			t.Fatalf("could not stop redis %s", err)
		}
	}()

	return redis.NewClient(&redis.Options{Addr: "127.0.0.1:4420"})
}

func TestAdd(t *testing.T) {

	cases := []struct {
		name             string
		payload          ForecastPayload
		key              string
		expectedResponse ForecastPayload
		err              string
	}{
		{
			// not really valid is it if half the fields are missing
			name: "Test Valid Add",
			payload: ForecastPayload{
				LastModelRunTime: "10",
				ForecastTotals:   map[string]*ForecastTotals{"1": {HighestTemp: 10.00}},
				Windows: [][]time.Time{
					{
						time.Date(2022, 2, 2, 2, 2, 2, 0, time.UTC),
					},
				},
			},
			key: "key",
			expectedResponse: ForecastPayload{
				LastModelRunTime: "10",
				ForecastTotals:   map[string]*ForecastTotals{"1": {HighestTemp: 10.00}},
				Windows: [][]time.Time{
					{
						time.Date(2022, 2, 2, 2, 2, 2, 0, time.UTC),
					},
				},
			},
		},
	}

	rdb := NewClient(t)
	store := NewMetStore(rdb, NewLogger("testlog.txt"))

	for _, tc := range cases {
		if err := store.Add(context.Background(), tc.key, tc.payload); err != nil {
			t.Fatalf("failed storing %s", err)
		}
		res, err := store.Rdb.Get(context.Background(), tc.key).Bytes()
		if err != nil {
			t.Fatalf("Failed to get from redis %s", err)
		}

		var testRes *ForecastPayload

		if err := json.Unmarshal(res, &testRes); err != nil {
			t.Fatalf("failed decoding test response %s", err)
		}

		assert.Equal(t, tc.expectedResponse, *testRes)
	}

}

// cases := []struct {
// 	name        string
// 	expectedVal map[string]*ForecastTotals
// }{
// 	{
// 		name:        "valid",
// 		expectedVal: map[string]*ForecastTotals{"1": &ForecastTotals{HighestTemp: 10.00}},
// 	},
// }
