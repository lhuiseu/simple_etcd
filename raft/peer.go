package raft

type Peer struct {
	ID uint64
	Context []byte
}
