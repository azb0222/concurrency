package mutex

import (
	"fmt"
	"sync"
	"time"
)

func Main() {
	MutexEx()
	RWMutexEx()
}

func MutexEx() {

	var criticalSection int
	var lock sync.Mutex
	var wg sync.WaitGroup

	increment := func() {
		lock.Lock()
		defer lock.Unlock() //will execute in case of a panic to avoid deadlock
		criticalSection++
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		criticalSection--
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decrement()
		}()
	}

	wg.Wait()
	fmt.Printf("Arithmetic complete %v", criticalSection)
}

/*
sync.RWMutex allows multiple concurrent processes to read data at the same time unless lock is being held for writing
use sync.RWMutex whenever possible to limit time spent in critical sections
*/
func RWMutexEx() {
	criticalSection := 0
	var lock sync.RWMutex
	rLock := lock.RLocker() //lock for reading; will be granted read access to data unless lock is being held for writing; therefore multiple observers can read at the same time if nothing is being produced
	var wg sync.WaitGroup

	slowProducer := func(wg *sync.WaitGroup, l sync.Locker, data *int) {
		defer wg.Done()
		l.Lock()
		*data += 1
		defer l.Unlock()
		time.Sleep(2) //to simulate the producer goroutine being slower than the observer goroutines: "taking time to produce data"
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker, data *int) {
		defer wg.Done()
		l.Lock()
		fmt.Printf("Data is: %v", *data)
		defer l.Unlock()
	}

	producerNum := 5
	observerNum := 5

	wg.Add(producerNum + observerNum)

	for i := 0; i < producerNum; i++ {
		go slowProducer(&wg, rLock, &criticalSection)
	}

	for i := 0; i < observerNum; i++ {
		go observer(&wg, rLock, &criticalSection)
	}

	wg.Wait()
}
