package brpc

// NewServer creates a new server according to the protocol.
func NewServer(protocol string, options ...ServerOption) Server {
	proto := getProtocol(protocol)
	if proto == nil {
		panic(ErrProtocolNotFound)
	}
	return proto.NewServer(options...)
}
