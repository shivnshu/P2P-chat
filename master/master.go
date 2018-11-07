package master

import (
	"log"
	"net/http"

	"github.com/shivnshu/P2P-chat/common/iface"
)

type PeerInfo struct {
	IP   string
	Port int
}

type Master struct {
	PeersInfo []PeerInfo
}

func (c Master) startMasterNode(args iface.CommonArgs) {
	master_ip := args.IP
	master_port := args.Port
	err := c.startListening(master_ip, master_port)
	if err != nil {
		log.Fatalf("Unable to start listening @%s:%d", master_ip, master_port)
	}
}

func (c Master) startListening(ip string, port int) error {
	addr := iface.GetAddress(ip, port).String()
	log.Println("Listening on", addr)
	http.HandleFunc("/", c.requestHandler)
	err := http.ListenAndServe(addr, nil)
	return err
}

func (c Master) requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got a request from %s", r.RemoteAddr)
}
