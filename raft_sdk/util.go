package raft_sdk

import(
	"math/rand"
	"sync"
	"time"
)
type lockedRand struct{
	mu sync.Mutex
	rand *rand.Rand
}

func (lock *lockedRand) Intn(n int) int{
	lock.mu.Lock()
	v := lock.rand.Intn(n)
	lock.mu.Unlock()

	return v

}

var globalRand = &lockedRand{
	rand : rand.New(rand.NewSource(time.Now().UnixNano())),
}