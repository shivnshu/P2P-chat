package peer

import (
	"bufio"
	"fmt"
	"github.com/shivnshu/P2P-chat/common/iface"
	"net/http"
	"os"
	"strconv"
)

type Peer struct {
	Self       iface.PeerInfo
	Neighbours []iface.PeerInfo
}

func (c Peer) startPeerNode(args iface.CommonArgs) {
	c.Self.Alias = args.Alias

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter master node IP: ")
	master_ip, _ := reader.ReadString('\n')
	master_ip = master_ip[:len(master_ip)-1]
	fmt.Print("Enter master node port: ")
	var master_port int
	_, err := fmt.Scanf("%d", &master_port)

	// fmt.Println(master_ip, master_port)
	err = c.registerWithMaster(master_ip, master_port)
	if err != nil {
		panic(err)
	}
}

func (c Peer) registerWithMaster(ip string, port int) error {
	master_url := "http://" + ip + ":" + strconv.Itoa(port)
	fmt.Println("hello", master_url)
	req, err := http.NewRequest("POST", master_url, nil)
	fmt.Println(req, err)
	return nil
}
