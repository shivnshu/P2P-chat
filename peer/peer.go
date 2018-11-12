package peer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/marcusolsson/tui-go"
	"github.com/shivnshu/P2P-chat/common/iface"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Peer struct {
	Self               iface.PeerInfo
	Neighbours         []iface.PeerInfo
	NeighboursLock     sync.Mutex
	ChatHistoryBox     *tui.Box
	ChatHistoryBoxLock sync.Mutex
	UIPainter          tui.UI
	ReadMsgs           map[string]bool
	ReadMsgsLock       sync.Mutex
}

// Entry point for peer node
func (c *Peer) startPeerNode(args iface.CommonArgs) {
	c.Self.IP = args.IP
	c.Self.Port = args.Port
	c.Self.Alias = args.Alias

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter master node IP: ")
	master_ip, _ := reader.ReadString('\n')
	master_ip = master_ip[:len(master_ip)-1]
	fmt.Print("Enter master node port: ")
	var master_port int
	_, err := fmt.Scanf("%d", &master_port)

	go c.periodicMasterRegistration(master_ip, master_port)

	// Initialize ReadMsgs map
	c.ReadMsgs = make(map[string]bool)
	go c.startChatBox()

	err = c.startListening()
	if err != nil {
		fmt.Println("Unable to start listening")
		return
	}
}

func (c *Peer) periodicMasterRegistration(ip string, port int) {
	for {
		neighbours, err := c.registerWithMaster(ip, port)
		if err != nil {
			// panic(err)
			time.Sleep(iface.MasterRegistrationInterval * time.Second)
			continue
		}
		c.NeighboursLock.Lock()
		c.Neighbours = neighbours
		// c.printToChat(strconv.Itoa(len(c.Neighbours)))
		c.NeighboursLock.Unlock()
		time.Sleep(iface.MasterRegistrationInterval * time.Second)
	}
}

// Send request to master node
func (c *Peer) registerWithMaster(ip string, port int) ([]iface.PeerInfo, error) {
	master_url := "http://" + ip + ":" + strconv.Itoa(port)
	master_url = master_url + "?buffer_size=" + strconv.Itoa(iface.DefaultBufferSize)
	jsonBytes, _ := json.Marshal(c.Self)
	req, err := http.NewRequest("POST", master_url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// panic(err)
		return nil, err
	}
	defer resp.Body.Close()

	var neighbours []iface.PeerInfo
	err = json.NewDecoder(resp.Body).Decode(&neighbours)
	if err != nil {
		return nil, err
	}
	return neighbours, nil
}

// Start listening for messages on user supplied port
func (c *Peer) startListening() error {
	addr := iface.GetAddress(c.Self.IP, c.Self.Port).String()
	fmt.Println("Listening on", addr)
	http.HandleFunc("/", c.msgHandler)
	err := http.ListenAndServe(addr, nil)
	return err
}

// Handler for any message arrival to its port
// Accept if it is destined for it else send to all its neighbours
func (c *Peer) msgHandler(w http.ResponseWriter, r *http.Request) {
	var msg iface.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		return
	}
	if msg.TTL <= 0 {
		return
	}
	// Async as to escape from deadlock
	go c.asyncMsgHandler(msg)
}

func (c *Peer) asyncMsgHandler(msg iface.Message) {
	msg.TTL--
	if msg.ToAlias == c.Self.Alias {
		c.recvMessage(msg)
	} else {
		if msg.ToAlias == "ALL" {
			c.recvMessage(msg)
		}
		c.sendMessage(msg)
	}
}

// Send msg to all its neighbours
func (c *Peer) sendMessage(msg iface.Message) {
	client := &http.Client{}

	var url string

	c.NeighboursLock.Lock()
	// c.printToChat(strconv.Itoa(len(c.Neighbours)))
	for _, peer := range c.Neighbours {
		url = "http://" + peer.IP + ":" + strconv.Itoa(peer.Port)
		jsonBytes, _ := json.Marshal(msg)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
	}
	c.NeighboursLock.Unlock()
}
