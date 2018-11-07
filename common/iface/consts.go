package iface

const (
	DefaultPort = 2222
)

type CommonArgs struct {
	Type  string
	IP    string
	Port  int
	Alias string
}
