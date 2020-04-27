package nodelock

import "sync"

type Locker struct {
	locks chan int
}

func (l *Locker) Lock() {
	l.locks <- 0
}

func (l *Locker) Unlock() {
	<-l.locks
}

func NewLocker(permits uint) sync.Locker {
	return &Locker{locks: make(chan int, permits)}
}
