package master

import (
	"encoding/json"
	"fmt"
	"github.com/shivnshu/P2P-chat/common/iface"
	"net/http"
	"strconv"
	"sync"
	"time"
    "math/rand"
)

type Master struct {
	PeersInfo     []iface.PeerInfo
	PeersInfoLock sync.Mutex
}

func (c *Master) startMasterNode(args iface.CommonArgs) {
	master_ip := args.IP
	master_port := args.Port
	go c.peersInfoPolling()
	err := c.startListening(master_ip, master_port)
	if err != nil {
		fmt.Printf("Unable to start listening @%s:%d\n", master_ip, master_port)
		return
	}
}

func (c *Master) peersInfoPolling() {
	for {
		now := time.Now()
		secs := now.Unix()
		tmpPeersInfo := []iface.PeerInfo{}
		c.PeersInfoLock.Lock()
		for _, val := range c.PeersInfo {
			if val.TimeStamp > secs-iface.MasterPeerInfoTimeout {
				tmpPeersInfo = append(tmpPeersInfo, val)
			}
		}
		c.PeersInfo = tmpPeersInfo
		c.PeersInfoLock.Unlock()
		time.Sleep(iface.MasterPeerInfoTimeout * time.Second)
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

func (c *Master) getPeersInfo(buffer_size int) []iface.PeerInfo {
	var result []iface.PeerInfo
    r := rand.New(rand.NewSource(time.Now().Unix()))
	c.PeersInfoLock.Lock()
    for ctr, i := range r.Perm(len(c.PeersInfo)) {
        if ctr >= buffer_size {
            break
        }
        result = append(result, c.PeersInfo[i]);
    }
	c.PeersInfoLock.Unlock()
	return result
}

func (c *Master) addToPeersInfo(peerInfo iface.PeerInfo) {
	if c.existsInPeersInfo(peerInfo) == false {
		c.PeersInfoLock.Lock()
		now := time.Now()
		peerInfo.TimeStamp = now.Unix()
		c.PeersInfo = append(c.PeersInfo, peerInfo)
		c.PeersInfoLock.Unlock()
	}
	return
}

func (c *Master) existsInPeersInfo(peerInfo iface.PeerInfo) bool {
	res := false
	c.PeersInfoLock.Lock()
	for i, val := range c.PeersInfo {
		if val.Alias == peerInfo.Alias {
			res = true
			// Update timestamp
			now := time.Now()
			c.PeersInfo[i].TimeStamp = now.Unix()
			break
		}
	}
	c.PeersInfoLock.Unlock()
	return res
}
