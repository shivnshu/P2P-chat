package peer

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/shivnshu/P2P-chat/common/iface"
)

func Init() error {
	var peer_args iface.CommonArgs
	arg.MustParse(&peer_args)

	if peer_args.IP == "" || peer_args.Port == 0 {
		return fmt.Errorf("Please provide master ip address and port number.")
	}

	peer := Peer{}
	peer.startPeerNode(peer_args)

	return nil
}