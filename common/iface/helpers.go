package iface

import (
	"strconv"
)

type Address struct {
	IP   string
	Port int
}

func GetAddress(ip string, port int) Address {
	return Address{
		IP:   ip,
		Port: port,
	}
}

func (a Address) String() string {
	return a.IP + ":" + strconv.Itoa(a.Port)
}
