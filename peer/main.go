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

	if peer_args.Alias == "" {
		return fmt.Errorf("Please provide alias to be used to recognize the peer.")
	}

	if peer_args.BufferSize == 0 {
		fmt.Printf("Buffer size not provided, using %d as buffer size.\n", iface.DefaultBufferSize)
		peer_args.BufferSize = iface.DefaultBufferSize
	}

	peer := Peer{}
	peer.startPeerNode(peer_args)

	return nil
}
