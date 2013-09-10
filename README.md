go-etcd-lock
============

A simple lock library on top of etcd in Go.

Locks are acquired with a timeout which will be automatically refreshed so long as the node with the lock remains alive.  If a node holding a lock were to die, the lock would automatically
release after the timeout period, and a new node would be granted the lock.

# Example

    import "github.com/hayesgm/go-etcd-lock/lock"

    goChan, stopChan := lock.Acquire(cli, "mylock", 20) // Try to get a lock with 20 second timeout

    go func() {
      <- goChan // Wait to Acquire lock
      
      for run := true; run; {
        select {
        case <-stopChan:
          run = false // We're going to exit
        default:
          log.Println("I am king")
        }
      }
    }

