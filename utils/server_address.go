package utils

import "strings"

type ServerAddress struct {
	protocol string
	host     string
	port     string
}

func NewServerAddress(protocol string, host string, port string) *ServerAddress {
	return &ServerAddress{protocol: protocol, host: host, port: port}
}

func (s *ServerAddress) String() string {
	var serverAddress strings.Builder
	serverAddress.WriteString(s.host)
	serverAddress.WriteString(":")
	serverAddress.WriteString(s.port)

	return serverAddress.String()
}

func (s *ServerAddress) StringInclProtocol() string {
	var serverAddress strings.Builder
	serverAddress.WriteString(s.protocol)
	serverAddress.WriteString("://")
	serverAddress.WriteString(s.host)
	serverAddress.WriteString(":")
	serverAddress.WriteString(s.port)

	return serverAddress.String()
}
