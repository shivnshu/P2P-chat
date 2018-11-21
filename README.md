# P2P-chat

## Directory Structure
```
.
├── README.md
├── main.go ................: Application entrypoint
├── common
│   └── iface
│       ├── consts.go ......: Common structs and default values
│       └── helpers.go .....: Helper functions
├── master
│   ├── main.go ............: Master node entrypoint
│   └── master.go ..........: Master node methods
└── peer
    ├── main.go ............: Peer node entrypoint
    ├── peer.go ............: Peer node methods
    └── chatbox.go .........: Terminal UI based on curses
```
## Usage
* Run `go run main.go --type=master --port=<port_num>` to run the master node.
* Run `go run main.go --type=peer --ip=<public_ip> --port=<port_num> --alias=<alias>` to run a peer node.
