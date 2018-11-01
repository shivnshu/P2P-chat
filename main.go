package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/shivnshu/P2P-chat/common/iface"
	"github.com/shivnshu/P2P-chat/master"
	// "github.com/shivnshu/P2P-chat/client"
)

func main() {
	var runType iface.CommonArgs
	arg.Parse(&runType)

	var err error

	switch runType.Type {
	case "master":
		err = master.Init()
	case "client":
	default:
		fmt.Println("Specify type for the node - master or client")
	}

	if err != nil {
		log.Fatal(err)
	}
}
