package main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func startRedis() (*exec.Cmd, error) {
	cmd := exec.Command("redis-server", "--port", "4420")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	time.Sleep(2 * time.Second)
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

	return redis.NewClient(&redis.Options{Addr: "localhost:4420"})
}

func TestAdd(t *testing.T) {

	rdb := NewClient()

}
