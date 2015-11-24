redislock
=========

Pessimistic locking for Go using Redis.

### Installation

    go get -u github.com/everalbum/redislock

### Documentation

http://godoc.org/github.com/everalbum/redislock

### Example

```go
m := redislock.NewMutex(conn, "user:123")
// optionally set timeout with: m.Timeout = time.Minute
ok, err := m.TryLock()
if err != nil {
	log.Fatal("Error while attempting lock")
}
if !ok {
	// User is in use - return to avoid duplicate work, race conditions, etc.
	return
}
defer lock.Unlock()

// Do something with the user.
```
