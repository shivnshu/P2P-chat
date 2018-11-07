package peer

import (
	"github.com/shivnshu/P2P-chat/common/iface"
)

type Peer struct {
	Alias      string
	Neighbours []*Peer
}

func (c Peer) startPeerNode(args iface.CommonArgs) {
	master_ip := args.IP
	master_port := args.Port
	c.Alias = args.Alias

	err := c.registerWithMaster(master_ip, master_port)
	if err != nil {
		panic(err)
	}
}

func (c Peer) registerWithMaster(ip string, port int) error {
	return nil
}
