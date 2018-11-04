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
	state StateType //raft的状态类型
	Vote uint64  //投票对象的id
	prs map[uint64]*Progress
	learnerPrs map[uint64]*Progress

	msgs []Message //发送给同伴的消息
}

type Config struct{
	ID uint64
	ElectionTick int
}

const (
	StateFollow StateType = iota
	StateCandidate
)

type StateType uint64

type stepFunc func (r *raft) error

func newRaft(c Config) *raft {

	return &raft{
		id : c.ID,
		electionTimeout: c.ElectionTick,
		prs: make(map[uint64]*Progress),
		learnerPrs: make(map[uint64]*Progress),
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

func (r *raft) becomeCandidater() {
	r.step = stepCandidate
	r.state = StateCandidate
	r.tick = r.tickElection
	r.reset(r.Term+1)
	r.Vote = r.id
}


func (r *raft) tickElection() {
	r.electionElapsed++
	if r.pastElectionTimeOut() {
		r.electionElapsed = 0
		//开始选举
		fmt.Println("election is begin.... , prepare send msg")
		//发送选举消息
		r.Step(Message{From:r.id, Type:MsgHup})
	}
}

func (r *raft) Step(m Message) error{
	switch {
	case m.Term == 0:
		// local message ,没有设置message的Term
	}

	switch m.Type {
	case MsgHup:
		fmt.Println(r.id, "is starting a new election at term ", r.Term)
		r.Campaign()
	default:
		fmt.Println("default")
	}
	return nil
}

/**
这个函数的名字有些费解
 */
func (r *raft) Campaign() {
	var voteMsg MessageType
	r.becomeCandidater()

	voteMsg = MsgVote
	for id := range r.prs {
		if id == r.id {
			continue
		}
		var ctx []byte
		fmt.Println(r.id," begin to send msg to peers:",id," at term:", r.Term)
		r.send(Message{Term:r.Term, To:id, Type:voteMsg, Context:ctx,})
	}
}



/**
判断是否已经超过选举的时机
 */
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

/**
成为新角色后的重置
 */
func (r *raft) reset(term uint64) {
	fmt.Println("origin term：",r.Term)
	if (r.Term != term) {
		r.Term = term
	}
	r.electionElapsed = 0
	r.resetRandomizedElectionTimeout()
}


func stepFollow(r *raft) error {

	return nil
}

func stepCandidate(r *raft) error{
	return nil
}

/**
产生随机选举时间
 */
func (r *raft) resetRandomizedElectionTimeout() {
	//为了防止各个node之前出现candidate选举竞态，每个node的选举超时时间都是随机的
	r.randomizedElectionTimeout = r.electionTimeout + globalRand.Intn(r.electionTimeout)
	fmt.Println("rand time:",r.randomizedElectionTimeout)

}


/**
生成Progress相关的
 */
func (r *raft) addNode(id uint64) {
	r.addNodeOrLearnerNode(id, false)
}

func (r *raft) addNodeOrLearnerNode(id uint64, isLearner bool) {
	ps := r.getProgress(id)
	if ps == nil {
		r.setProgress(id , isLearner)
	}

}

func (r *raft) getProgress(id uint64) *Progress{
	if pr,ok := r.prs[id]; ok {
		return pr
	}
	return r.learnerPrs[id]
}

func (r *raft) setProgress(id uint64, isLearner bool) {
	if !isLearner {
		delete(r.learnerPrs, id)
		r.prs[id] = Progress{}
		return
	}
}

/**
raft向同伴发送消息
 */

func (r *raft) send(msg Message) {
	msg.From = r.id
	r.msgs = append(r.msgs, msg)
	//明天继续看，后续是怎么实现的

}