package master

import (
	"github.com/alexflint/go-arg"
	"github.com/shivnshu/P2P-chat/common/iface"
	"log"
)

func Init() error {
	var master_args iface.CommonArgs
	arg.MustParse(&master_args)

	if master_args.Port == 0 {
		log.Printf("port not supplied, using default port %d", iface.DefaultPort)
		master_args.Port = iface.DefaultPort
	}

	if master_args.IP == "" {
		master_args.IP = "0.0.0.0"
	}

	master := Master{}
	master.startMasterNode(master_args)

	return nil
}
