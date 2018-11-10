package iface

const (
	DefaultPort       = 2222
	DefaultBufferSize = 2
)

type CommonArgs struct {
	Type       string
	IP         string
	Port       int
	Alias      string
	BufferSize int
}

type PeerInfo struct {
	IP    string
	Port  int
	Alias string
}
