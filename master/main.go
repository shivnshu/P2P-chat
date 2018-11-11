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
		fmt.Printf("port not supplied, using default port %d\n", iface.DefaultPort)
		master_args.Port = iface.DefaultPort
	}

	if master_args.IP == "" {
		master_args.IP = "0.0.0.0"
	}

	master := Master{}
	master.startMasterNode(master_args)

	return nil
}
