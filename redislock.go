/*
Package redislock implements a pessimistic lock using Redis.

For example, lock and unlock a user using its ID as a resource identifier:
	lock, ok, err := redislock.TryLock(conn, "user:123")
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		// User is in use - return to avoid duplicate work, race conditions, etc.
		return
	}
	defer lock.Unlock()

	// Do something with the user.
*/
package redislock

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
)

const DefaultTimeout = 10 * time.Minute

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

// Mutex represents a mutual exclusion lock.
type Mutex struct {
	Timeout  time.Duration // default value of DefaultTimeout
	resource string
	token    string
	conn     redis.Conn
}

func NewMutex(conn redis.Conn, resource string) *Mutex {
	return &Mutex{DefaultTimeout, resource, uuid.New(), conn}
}

// TryLock attempts to acquire a lock on the given resource in a non-blocking manner.
func (m *Mutex) TryLock() (ok bool, err error) {
	status, err := redis.String(m.conn.Do("SET", m.key(), m.token, "EX", int64(m.timeout()/time.Second), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return status == "OK", nil
}

// Unlock releases the lock. If the lock has timed out, it silently fails without error.
func (m *Mutex) Unlock() (err error) {
	_, err = unlockScript.Do(m.conn, m.key(), m.token)
	return
}

func (m *Mutex) key() string {
	return fmt.Sprintf("redislock:%s", m.resource)
}

func (m *Mutex) timeout() time.Duration {
	if m.Timeout == 0 {
		return DefaultTimeout
	}
	return m.Timeout
}
