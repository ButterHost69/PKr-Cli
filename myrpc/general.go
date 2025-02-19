package myrpc

import (
	"net/rpc"

	"github.com/ButterHost69/kcp-go"
)

const (
	SERVER_HANDLER_NAME = "Handler"
	CLIENT_HANDLER_NAME = "Handler"
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
