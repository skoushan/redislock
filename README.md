redislock
=========

Pessimistic locking for Go using Redis.

### Installation

    go get -u github.com/atomic-labs/redislock

### Documentation

http://godoc.org/github.com/atomic-labs/redislock

### Example

```go
lock, ok, err := redislock.TryLock(conn, "user:123")

if !ok {
	// User is in use - return to avoid duplicate work, race conditions, etc.
	return
}

if err != nil {
  log.Fatal("Error while attempting lock")
}

// Do something with the user.

// When you're done, unlock the user.
lock.Unlock()
```
