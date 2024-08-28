package signals

import (
	"fmt"
	"sync"
	"time"
)

func Main() {
	CondSignalEx()
}

// event: signal between 2(+) goroutines saying "something" has occurred
func CondSignalEx() {
	cond := sync.NewCond(&sync.Mutex{})           //NewCond() takes in a type that satisfies sync.Locker interface as a parameter
	criticalSection := make([]interface{}, 0, 10) //length = 0; capacity = 10

	slowRemoveFromQueue := func(delay time.Duration) {
		time.Sleep(delay)

		cond.L.Lock()
		criticalSection = criticalSection[1:] //remote 0th element in queue
		cond.L.Unlock()

		cond.Signal() //Signal() lets the longest waiting goroutine on the condition know that "something" has occurred; Go runtime maintains internal FIFO list of goroutines waiting to be signaled
	}

	for i := 0; i < 10; i++ {
		cond.L.Lock()

		for len(criticalSection) == 2 { //will suspend execution for as long as queue length == 2
			cond.Wait() //suspend current goroutine until we receive notification that something has occurred from Signal()
			//upon entering cond.Wait, Unlock is called on cond.L; upon exiting cond.Wait, Lock is called on cond.L
		}

		criticalSection = append(criticalSection, struct{}{})
		go slowRemoveFromQueue(2 * time.Second)
		cond.L.Unlock()
	}
}

func CondBroadcastEx() {
	//unlike Signal(), Broadcast will send a signal to ALL waiting goroutines

	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			wg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait() //suspends execution until awoken by Broadcast or Singal
			fn()
		}()

		wg.Wait()
	}
	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		time.Sleep(1)
		fmt.Println("Maximizing Window.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		time.Sleep(1)
		fmt.Println("Displaying a Popup.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		time.Sleep(1)
		fmt.Println("Mouse Clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast() //trigger all three subscriber functions

	clickRegistered.Wait() //do not exit main goroutine until all three subscriber functions finish executing
}
