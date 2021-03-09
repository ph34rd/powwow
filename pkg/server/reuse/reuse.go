package reuse

import (
	"context"
	"net"
)

var conf = net.ListenConfig{
	Control: Control,
}

// Listen announces on the local network address.
// Additional options are set SO_REUSEPORT or SO_REUSEADDR.
func Listen(network, address string) (net.Listener, error) {
	return conf.Listen(context.Background(), network, address)
}
