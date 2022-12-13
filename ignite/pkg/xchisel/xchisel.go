package xchisel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	chclient "github.com/jpillora/chisel/client"
	chserver "github.com/jpillora/chisel/server"
)

var DefaultServerPort = "7575"

func ServerAddr() string {
	return os.Getenv("CHISEL_ADDR")
}

func IsEnabled() bool {
	return ServerAddr() != ""
}

func StartServer(ctx context.Context, port string) error {
	s, err := chserver.NewServer(&chserver.Config{})
	if err != nil {
		return err
	}
	if err := s.StartContext(ctx, "127.0.0.1", port); err != nil {
		return err
	}
	if err = s.Wait(); errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func StartClient(ctx context.Context, serverAddr, localPort, remotePort string) error {
	c, err := chclient.NewClient(&chclient.Config{
		MaxRetryInterval: time.Second,
		MaxRetryCount:    -1,
		Server:           serverAddr,
		Remotes:          []string{fmt.Sprintf("127.0.0.1:%s:127.0.0.1:%s", localPort, remotePort)},
	})
	if err != nil {
		return err
	}
	c.Logger.Info = false
	c.Logger.Debug = false
	if err := c.Start(ctx); err != nil {
		return err
	}
	if err = c.Wait(); errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}
