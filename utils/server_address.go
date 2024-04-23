package utils

import "strings"

var serverAddressSeparator = ":"

type ServerAddress struct {
	host string
	port string
}

func NewServerAddress(host string, port string) *ServerAddress {
	return &ServerAddress{host: host, port: port}
}

func (s *ServerAddress) String() string {
	var serverAddress strings.Builder
	serverAddress.WriteString(s.host)
	serverAddress.WriteString(serverAddressSeparator)
	serverAddress.WriteString(s.port)

	return serverAddress.String()
}
