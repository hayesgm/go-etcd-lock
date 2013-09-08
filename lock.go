package lock

import (
  "log"
  "github.com/coreos/go-etcd/etcd"
  "github.com/coreos/etcd/store"
  "github.com/go-contrib/uuid"
  "time"
)

// This will attempt to Acquire a lock for a given amount of time.
// If the lock is acquired, goChan will be unblocked
// If the lock is lost, stopChan will be unblocked
func Acquire(cli *etcd.Client, lock string, timeout uint64) (goChan chan int, stopChan chan int) {
  me := uuid.NewV1().String()
  goChan = make(chan int)
  stopChan = make(chan int)

  go func() {
    // log.Println("Hello, I am:",me)

    for {
      resp, acq, err := cli.TestAndSet(lock, "", me, timeout)

      // log.Println("Lock Resp:",acq,resp,err)

      if !acq {
        // We want to watch for a change in the lock, and we'll repeat
        var watcherCh = make(chan *store.Response)
        var endCh = make(chan bool)

        go cli.Watch(lock, 0, watcherCh, endCh)
        <- watcherCh

        // Now, we'll try to acquire the lock, again
      } else {
        
        // We got a lock, we want to keep it
        go func() {
          for {
            resp, acq, err := cli.TestAndSet("/fiddler/watcher", me, me, timeout) // Keep the lock alive
            // log.Println("Reset Resp:",acq,resp,err)
            
            if !acq {
              // log.Println("Demoted:",me)
              stopChan <- 1 // Let's boot ourselves, we're no longer the leader
            }
            
            time.Sleep(time.Duration(timeout * 500) * time.Millisecond) // We'll re-up after 50% fo the lock period
          }
        }()
        // log.Println("King:",me)
        goChan <- 1
      }
    }
  }()

  return
}