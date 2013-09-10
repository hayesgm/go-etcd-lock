package daemon

import (
  "github.com/coreos/go-etcd/etcd"
  "github.com/hayesgm/go-etcd-lock/lock"
  "log"
)

/*
  Daemon is a simple routine that will run on one go server
  in an etcd cluster.
*/

// DaemonFunc takes a stopCh and should stop when this is unblocked
type DaemonFunc func(stopCh chan int)

// Ensures one daemon is running through a distributed lock on etcd
// This is a non-blocking call
// TODO: Handler errors
func RunOne(cli *etcd.Client, name string, f DaemonFunc, timeout int) {
  goChan, stopChan := lock.Acquire(cli, name, uint64(timeout)) // Try to get a lock with 20 second timeout

  go func() {
    log.Println("Waiting patiently")
    <- goChan // Wait to Acquire lock
    log.Println("Running daemon")
    f(stopChan)
  }()
}