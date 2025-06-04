package myrpc

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/ButterHost69/kcp-go"
)

const (
	SERVER_HANDLER_NAME      = "Handler"
	CLIENT_BASE_HANDLER_NAME = "ClientHandler"
)

func call(rpcname string, args interface{}, reply interface{}, ripaddr, lipaddr string) error {
	conn, err := kcp.Dial(ripaddr, lipaddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := rpc.NewClient(conn)
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err != nil {
		return err
	}

	return nil
}

func callWithContextAndConn(ctx context.Context, rpcname string, args interface{}, reply interface{}, ripaddr string, udpConn *net.UDPConn) error {
	// Dial the remote address
	kcpConn, err := kcp.DialWithConnAndOptions(ripaddr, nil, 0, 0, udpConn)
	if err != nil {
		return err
	}
	defer kcpConn.Close()
	kcpConn.SetNoDelay(0, 15000, 0, 0)
	kcpConn.SetDeadline(time.Now().Add(30 * time.Second)) // Overall timeout
	kcpConn.SetACKNoDelay(false)                          // Batch ACKs to reduce traffic

	// Find a Way to close the kcp conn without closing UDP Connection
	// defer conn.Close()

	c := rpc.NewClient(kcpConn)
	defer c.Close()

	// Create a channel to handle the RPC call with context
	done := make(chan error, 1)
	go func() {
		done <- c.Call(rpcname, args, reply)
	}()

	select {
	case <-ctx.Done():
		if err := c.Close(); err != nil {
			return fmt.Errorf("RPC call timed out - %s\nAlso Error in Closing RPC %v", ripaddr, err)
		}
		return fmt.Errorf("RPC call timed out - %s", ripaddr)
	case err := <-done:
		if cerr := c.Close(); err != nil {
			return fmt.Errorf("%v, Also Error in Closing RPC %v", err, cerr)
		}
		return err
	}
}
