package iface

import (
	"crypto/md5"
	"encoding/hex"
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

func CalculateMD5Hash(msg Message) string {
	sumStr := msg.ToAlias + ":" + msg.FromAlias + ":" + msg.Msg + ":" + msg.Time.String()
	hasher := md5.New()
	hasher.Write([]byte(sumStr))
	return hex.EncodeToString(hasher.Sum(nil))
}
