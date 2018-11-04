package raft_sdk

import "fmt"

type Node interface {
	Tick();
}

type node struct {
	tickc  chan struct{}

}

func StartNode(c Config) Node {
	fmt.Println("启动库node节点....")
	//new SDK raft. raft并未作为Node的属性而存在，而是在node run的时候，作为参数传入
	raft := newRaft(c)

	raft.becomeFollow(1)

	//new SDK Node
	n := newNode()

	//异步run Node
	go n.run(raft)

	return n
}

func newNode() *node {
	return &node {
		tickc : make(chan struct{}),
	}
}

func (n *node) run(raft *raft) {
	fmt.Println("node 库启动异步协程，监听时钟信号....")
	for{
		select {
		//Node收到定时chanel的消息。当node收到chanel信号的时候，所有的行为都作用在raft上面，这点需要注意下
		case <-n.tickc:
			raft.tick()

		}
	}
}

func (n *node) Tick() {
	//当到达时钟周期时（100ms）,raftNode的serveChannel方法会调用这个函数
	//注意，外层不是直接操作node的tickc管道，因为是私有变量
	select {
	case n.tickc <- struct{}{}: //空struct不占用存储空间
	default:
		fmt.Println("....")
	}
}