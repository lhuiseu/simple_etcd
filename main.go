package main

import(
	"flag"
	"github.com/simple_etcd/raft"
	"strings"
	"fmt"
)

func main() {
	cluster := flag.String("cluster","http://127.0.0.1:9999","comma separated cluster peers")
	id := flag.Int("ID", 1, "cluster id number")
	flag.Parse()

	errorC := make(chan error)

	raft.NewRaftNode(*id, strings.Split(*cluster, ","))
	fmt.Println("end   ",*id,  *cluster)

	if _, ok := <-errorC ; ok{
		fmt.Println("end")
	}
}


