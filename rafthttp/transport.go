package rafthttp

import (
	"net/http"
	"path"
)
type Transport struct {

}

var (
	RaftPrefix         = "/raft"
	RaftStreamPrefix   = path.Join(RaftPrefix, "stream")
)



func (transport *Transport) Handler() http.Handler{
	mux := http.NewServeMux()
	streamHandler := newStraemHandler()
	mux.Handle(RaftStreamPrefix+"/", streamHandler)
	return mux
}

//自带处理器
type streamHandler struct{

}
//处理器需要实现ServeHTTP的方法
func (streamHandler *streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}


//生成处理其的函数
func newStraemHandler() *streamHandler{
	return &streamHandler{

	}
}