package iface

import "time"

const (
	DefaultPort       = 2222
	DefaultBufferSize = 2
	DefaultTTL        = 3
	MAX_PEERS         = 10
)

type CommonArgs struct {
	Type       string
	IP         string
	Port       int
	Alias      string
	BufferSize int
}

type PeerInfo struct {
	IP    string
	Port  int
	Alias string
}

type Message struct {
	ToAlias   string
	FromAlias string
	Msg       string
	Time      time.Time
	MD5Hash   string
	TTL       int
}
