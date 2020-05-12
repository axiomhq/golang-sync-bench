package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type Taker func()
type Dropper func()

// Uncontested primitives
var ch0 = make(chan int)
var ch1 = make(chan int, 1)
var ch10 = make(chan int, 10)
var lock = &sync.Mutex{}
var rwlock = &sync.RWMutex{}
var atomicInt64 = int64(-1)

// Contested primitives
var ch0Spam = make(chan int)
var ch1Spam = make(chan int, 1)
var ch10Spam = make(chan int, 10)
var lockSpam = &sync.Mutex{}
var rwlockSpam = &sync.RWMutex{}
var atomicInt64Spam = int64(-1)

func checkChan(ch chan int) int {
	select {
	case val := <-ch:
		return val
	default:
		return 0
	}
}

func spamChan(ch chan int) {
	for i := 0; ; i++ {
		ch <- i
	}
}

var spamMutexResult = -1 // Take make sure stuff is not optimized away
func spamMutex(m *sync.Mutex) {
	spamMutexResult = 0
	for i := 0; ; i++ {
		m.Lock()
		spamMutexResult = i * i / 17
		m.Unlock()

		// ?! Below Gosched() call is required for this thread to ever relinquish
		// control of the mutex. Other threads waiting to lock it will *not* be woken up
		// without it. (deadlocking the test)
		// On the other hand calling out to the scheduler gives a massive hit on the
		// *entire* benchmark suite - not just the SpamMutex one!
		runtime.Gosched()
	}
}

var spamRWMutexResult = -1 // To make sure stuff is not optimized away
func spamRWMutex(m *sync.RWMutex) {
	spamRWMutexResult = 0
	for i := 0; ; i++ {
		m.Lock()
		spamRWMutexResult = i * i / 17
		m.Unlock()
	}
}

func spamAtomicInt64(cnt *int64) {
	for {
		atomic.AddInt64(cnt, 1)

		// ?! See note in spamMutex()
		runtime.Gosched()
	}
}

func DoWork(i int) int {
	mod := i % 7
	if mod < 0 {
		mod = -mod
	}
	return (i + 1) << uint(mod)
}

var notused int64 = 1

func BenchmarkSync(b *testing.B) {
	go spamChan(ch0Spam)
	go spamChan(ch1Spam)
	go spamChan(ch10Spam)
	//go spamMutex(lockSpam) // disabled -- see comment in spamMutex()
	go spamRWMutex(rwlockSpam)
	//go spamAtomicInt64(&atomicInt64Spam) // disabled -- see comment in spamAtomicInt64()

	var tests = []struct {
		Name string
		Take Taker
		Drop Dropper
	}{
		{"Atomic", func() { atomic.LoadInt64(&atomicInt64) }, func() { atomic.AddInt64(&atomicInt64, 1) }},
		{"ch0", func() { checkChan(ch0) }, func() {}},
		{"ch1", func() { checkChan(ch1) }, func() {}},
		{"ch10", func() { checkChan(ch10) }, func() {}},
		{"Mutex", func() { lock.Lock() }, func() { lock.Unlock() }},
		{"RWMutex", func() { rwlock.RLock() }, func() { rwlock.RUnlock() }},
		//{"SpamAtomic", func(){atomic.LoadInt64(&atomicInt64Spam)}, func(){atomic.AddInt64(&atomicInt64Spam, 1)}}, // disabled, see comment
		{"Spamch0", func() { checkChan(ch0Spam) }, func() {}},
		{"Spamch1", func() { checkChan(ch1Spam) }, func() {}},
		{"Spamch10", func() { checkChan(ch10Spam) }, func() {}},
		//{"SpamMutex", func(){lockSpam.Lock()}, func(){lockSpam.Unlock()}}, // disabled, see above
		{"SpamRWMutex", func() { rwlockSpam.Lock() }, func() { rwlockSpam.Unlock() }},
	}

	for _, test := range tests {
		b.Run(test.Name, func(b *testing.B) {
			j := 0
			for i := 0; i < b.N; i++ {
				test.Take()
				j = DoWork(j)
				test.Drop()
			}

			if j == 0 {
				b.Fatal("j == 0")
			}
		})

		// You're not gonna believe this:
		// If you comment out the seemingly useless if-statement below then the Atomic test runs 3x slower!
		if atomic.LoadInt64(&notused) == 0 {
			b.Fatal("Something is rotten")
		}
	}
}
