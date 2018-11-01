package raft

import(
	"net/url"
	"log"
	"net"
	"net/http"
	"github.com/simple_etcd/rafthttp"
	"fmt"
)

type raftNode struct {
	id int
	peers []string
	transport rafthttp.Transport
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

	rpeers := make([]Peer, len(rc.peers))
	for i := range rpeers {
		rpeers[i] = Peer{ID : uint64(i+1),}
	}

	go rc.serveRaft()


}

func (rc *raftNode) serveRaft() {

	url , err := url.Parse(rc.peers[rc.id-1])
	if (err != nil) {
		log.Fatalf("url is err (%v)", err)
	}
	fmt.Println("url:",url.Host)
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
