package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/shivnshu/P2P-chat/common/iface"
	"github.com/shivnshu/P2P-chat/master"
	"github.com/shivnshu/P2P-chat/peer"
)

func main() {
	var runType iface.CommonArgs
	arg.Parse(&runType)

	var err error

	switch runType.Type {
	case "master":
		err = master.Init()
	case "peer":
		err = peer.Init()
	default:
		fmt.Println("Specify type for the node - master or peer")
	}

	if err != nil {
		log.Fatal(err)
	}
}
