package network

import "net"

const (
	DEFAULT_IP   = "192.168.100.3"
	DEFAULT_PORT = "9000"
)

type Server struct {
	IP        string
	Port      string
	Clients   map[string]*net.Conn
	ConnCount int
}
