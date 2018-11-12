package master

import (
	"encoding/json"
	"fmt"
	"github.com/shivnshu/P2P-chat/common/iface"
	"net/http"
	"strconv"
	"sync"
)

type Master struct {
	PeersInfo     []iface.PeerInfo
	PeersInfoLock sync.Mutex
}

func (c *Master) startMasterNode(args iface.CommonArgs) {
	master_ip := args.IP
	master_port := args.Port
	err := c.startListening(master_ip, master_port)
	if err != nil {
		fmt.Printf("Unable to start listening @%s:%d\n", master_ip, master_port)
		return
	}
}

func (c *Master) startListening(ip string, port int) error {
	addr := iface.GetAddress(ip, port).String()
	fmt.Println("Listening on", addr)
	http.HandleFunc("/", c.requestHandler)
	err := http.ListenAndServe(addr, nil)
	return err
}

func (c *Master) requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a request from %s\n", r.RemoteAddr)
	buffer_size_str, present := r.URL.Query()["buffer_size"]
	if present != true {
		fmt.Println("GET paramater buffer_size not present")
		return
	}
	buffer_size, err := strconv.Atoi(buffer_size_str[0])
	if err != nil {
		fmt.Println("Unable to get buffer_size")
		return
	}
	var peerInfo iface.PeerInfo
	err = json.NewDecoder(r.Body).Decode(&peerInfo)
	if err != nil {
		fmt.Println("Unable to get peerInfo")
		return
	}

	peersInfo := c.getPeersInfo(buffer_size)
	c.addToPeersInfo(peerInfo)

	data, _ := json.Marshal(peersInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	fmt.Println("Response: ", peersInfo)
	return
}

// TODO: Randomly select at most buffer size no. of Peer Info and return it
func (c *Master) getPeersInfo(buffer_size int) []iface.PeerInfo {
	var result []iface.PeerInfo
	c.PeersInfoLock.Lock()
	result = c.PeersInfo
	c.PeersInfoLock.Unlock()
	return result
}

func (c *Master) addToPeersInfo(peerInfo iface.PeerInfo) {
	if c.existsInPeersInfo(peerInfo) == false {
		c.PeersInfoLock.Lock()
		c.PeersInfo = append(c.PeersInfo, peerInfo)
		c.PeersInfoLock.Unlock()
	}
	return
}

func (c *Master) existsInPeersInfo(peerInfo iface.PeerInfo) bool {
	res := false
	c.PeersInfoLock.Lock()
	for _, val := range c.PeersInfo {
		if val.Alias == peerInfo.Alias {
			res = true
			break
		}
	}
	c.PeersInfoLock.Unlock()
	return res
}
