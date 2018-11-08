package master

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/shivnshu/P2P-chat/common/iface"
)

type PeerInfo struct {
	IP    string
	Port  int
	Alias string
}

type Master struct {
	PeersInfo []PeerInfo
	size      int
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
	buffer_size_str, present := r.URL.Query()["buffer_size"]
	if present != true {
		log.Println("GET paramater buffer_size not present")
		return
	}
	buffer_size, err := strconv.Atoi(buffer_size_str[0])
	if err != nil {
		log.Println("Unable to get buffer_size")
		return
	}
	var peerInfo PeerInfo
	err = json.NewDecoder(r.Body).Decode(&peerInfo)
	if err != nil {
		log.Println("Unable to get peerInfo")
		return
	}

	peersInfo := c.getPeersInfo(buffer_size)
	c.addToPeersInfo(peerInfo)

	data, _ := json.Marshal(peersInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

// TODO: Take care of concurrency issues
func (c Master) getPeersInfo(buffer_size int) []PeerInfo {
	var result []PeerInfo
	return result
}

// TODO: Take care of concurrency issues
func (c Master) addToPeersInfo(peerInfo PeerInfo) {
	return
}
