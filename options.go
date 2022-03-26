package brpc

// DialOption configures how we set up the connection.
type DialOption interface{}

// CallOption configures how we call the server.
type CallOption interface{}

// ServerOption configures how we set up the server.
type ServerOption interface{}

// CallOptions is the options for calling.
type CallOptions struct {
}

// DialOptions is the options for dialing.
type DialOptions struct {
}

// ServerOptions is the options for new server.
type ServerOptions struct {
}
