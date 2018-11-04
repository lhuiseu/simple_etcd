package raft_sdk

type MessageType int32 //为什么是int32，而不是uint

const (
	MsgHup MessageType = 0
	MsgVote MessageType = 1
)


type Message struct{
	Type MessageType
	To uint64
	From uint64
	Term uint64
	Context []byte
}