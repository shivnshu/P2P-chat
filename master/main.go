package master

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/shivnshu/P2P-chat/common/iface"
)

func Init() error {
	var master_args iface.CommonArgs
	arg.MustParse(&master_args)

	if master_args.Port == 0 {
		return fmt.Errorf("Please provide port number.")
	}

	if master_args.IP == "" {
		master_args.IP = "0.0.0.0"
	}

	startMasterNode(master_args)

	return nil
}
