package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ph34rd/powwow/pkg/session"
)

const dialTimeout = 15 * time.Second

func RunClient(ctx context.Context, addr string) (string, error) {
	d := net.Dialer{Timeout: dialTimeout, Deadline: time.Now().Add(dialTimeout)}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return "", fmt.Errorf("dial err: %v", err)
	}
	defer conn.Close()

	c := session.NewClientSession(conn, mint)
	wow, err := c.GetWoW(ctx)
	if err != nil {
		return "", fmt.Errorf("get wow err: %v", err)
	}
	return wow, nil
}
