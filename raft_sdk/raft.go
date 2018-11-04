package raft_sdk

import "fmt"

type raft struct{
	id uint64
	Term uint64
	tick func()
	electionTimeout int
	electionElapsed int
	randomizedElectionTimeout int
	step stepFunc
}

type Config struct{
	ID uint64
	ElectionTick int
}

type stepFunc func (r *raft) error

func newRaft(c Config) *raft {

	return &raft{
		id : c.ID,
		electionTimeout: c.ElectionTick,
	}

}

func (r *raft) becomeFollow(term uint64) {
	fmt.Println("注册成为follow...")
	//成为follow后，step要做的事情
	r.step = stepFollow

	//reset 一些状态
	r.reset(term)
	//成为follow后，tick时钟到达后需要做的事情
	r.tick = r.tickElection
	fmt.Println("注册成为follow...end")

}


func (r *raft) tickElection() {
	r.electionElapsed++
	if r.pastElectionTimeOut() {
		r.electionElapsed = 0
		//开始选举
		fmt.Println("election is begin.... , prepare send msg")
	}
}

func (r *raft) pastElectionTimeOut() bool{

	//basd code write
	/**
	if (r.electionElapsed >= r.randomizedElectionTimeout) {
		return true
	}
	return false
	**/

	return r.electionElapsed >= r.randomizedElectionTimeout
}

func (r *raft) reset(term uint64) {
	if (r.Term != term) {
		r.Term = term
	}
	r.electionElapsed = 0
	r.resetRandomizedElectionTimeout()
}

func stepFollow(r *raft) error {

	return nil
}

func (r *raft) resetRandomizedElectionTimeout() {
	//为了防止各个node之前出现candidate选举竞态，每个node的选举超时时间都是随机的
	r.randomizedElectionTimeout = r.electionTimeout + globalRand.Intn(r.electionTimeout)
	fmt.Println("rand time:",r.randomizedElectionTimeout)

}