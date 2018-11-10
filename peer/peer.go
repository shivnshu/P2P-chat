package peer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shivnshu/P2P-chat/common/iface"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Peer struct {
	Self       iface.PeerInfo
	Neighbours []iface.PeerInfo
}

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

	neighbours, err := c.registerWithMaster(master_ip, master_port)
	if err != nil {
		panic(err)
	}

	c.Neighbours = neighbours

	err = c.startListening()
	if err != nil {
		fmt.Println("Unable to start listening")
		return
	}
}

func (c *Peer) registerWithMaster(ip string, port int) ([]iface.PeerInfo, error) {
	master_url := "http://" + ip + ":" + strconv.Itoa(port)
	master_url = master_url + "?buffer_size=" + strconv.Itoa(iface.DefaultBufferSize)
	jsonBytes, _ := json.Marshal(c.Self)
	req, err := http.NewRequest("POST", master_url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil, nil
}

func (c *Peer) startListening() error {
	addr := iface.GetAddress(c.Self.IP, c.Self.Port).String()
	fmt.Println("Listening on", addr)
	http.HandleFunc("/", c.msgHandler)
	err := http.ListenAndServe(addr, nil)
	return err
}

func (c *Peer) msgHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got a request from %s", r.RemoteAddr)
}
