package raft

import(
	"net/url"
	"log"
	"net"
	"net/http"
	"github.com/simple_etcd/rafthttp"
	"github.com/simple_etcd/raft_sdk"
	"fmt"
	"time"
)

type raftNode struct {
	id int
	peers []string
	transport rafthttp.Transport
	node raft_sdk.Node
	stopc chan struct{}
	httpstopc chan struct{}
	httpdonec chan struct{}

}

func NewRaftNode(id int, peers []string) {

	rc := &raftNode{
		id : id,
		peers : peers,
		transport : rafthttp.Transport{},
		stopc : make(chan struct{}),
		httpstopc : make(chan struct{}),
		httpdonec : make(chan struct{}),
	}
	go rc.startRaft()

}

func (rc *raftNode) startRaft() {

	rpeers := make([]raft_sdk.Peer, len(rc.peers))
	for i := range rpeers {
		rpeers[i] = raft_sdk.Peer{ID : uint64(i+1),}
	}

	//启动与之对应的Node...
	config := raft_sdk.Config {
		ID : uint64(rc.id), //为什么raftNode的id不是用的uint64呢？
		ElectionTick : 10, //当前raft节点的超时时间是10个周期,每个周期的时长，在serveChannel里面设置的是100ms
	}
	rc.node = raft_sdk.StartNode(config, rpeers)
	fmt.Println("start node lib OK...")

	go rc.serveRaft()
	go rc.serveChannel()


}

func (rc *raftNode) serveRaft() {

	url , err := url.Parse(rc.peers[rc.id-1])
	if (err != nil) {
		log.Fatalf("url is err (%v)", err)
	}
	//创建listener
	ln, err := net.Listen("tcp", url.Host)
	if (err != nil) {
		log.Fatalf("failed to listen (%v)", err)
	}

	//创建server
	err = (&http.Server{Handler: rc.transport.Handler()}).Serve(ln)
	fmt.Println("routine has not been block")
	select{
	case <-rc.stopc:
	default:
		fmt.Println("raftexample: Failed to serve rafthttp")
	}

	close(rc.httpdonec)

}

func (rc *raftNode) serveChannel() {

	//创建定时器
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rc.node.Tick()
		}
	}
}
