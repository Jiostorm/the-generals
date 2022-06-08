package network

import "net"

type Client struct {
	ID   string
	Conn net.Conn
}
