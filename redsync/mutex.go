package redsync

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/k81/kate/utils"
)

// Mutex is a distributed mutual exclusion lock.
type Mutex struct {
	name   string
	expiry time.Duration

	tries    int
	delayMin time.Duration
	delayMax time.Duration

	factor float64

	quorum int

	token string
	until time.Time

	nodem sync.Mutex

	pools []Pool
}

// GetToken return the token id for the mutex
func (m *Mutex) GetToken() string {
	return m.token
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	if m.token == "" {
		token, err := m.genToken()
		if err != nil {
			return err
		}
		m.token = token
	}

	for i := 0; i < m.tries; i++ {
		start := time.Now()

		n := 0
		for _, pool := range m.pools {
			ok := m.acquire(pool, m.token)
			if ok {
				n++
			}
		}

		until := time.Now().Add(m.expiry - time.Now().Sub(start) - time.Duration(int64(float64(m.expiry)*m.factor)) + 2*time.Millisecond)
		if n >= m.quorum && time.Now().Before(until) {
			m.until = until
			return nil
		}
		for _, pool := range m.pools {
			m.release(pool, m.token)
		}

		time.Sleep(m.getDelay())
	}

	return ErrFailed
}

// Unlock unlocks m and returns the status of unlock. It is a run-time error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	n := 0
	for _, pool := range m.pools {
		ok := m.release(pool, m.token)
		if ok {
			n++
		}
	}
	return n >= m.quorum
}

// Extend resets the mutex's expiry and returns the status of expiry extension. It is a run-time error if m is not locked on entry to Extend.
func (m *Mutex) Extend() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	n := 0
	for _, pool := range m.pools {
		ok := m.touch(pool, m.token, m.expiry)
		if ok {
			n++
		}
	}
	if n >= m.quorum {
		return true
	}

	for _, pool := range m.pools {
		m.release(pool, m.token)
	}
	return false
}

func (m *Mutex) genToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s_%s", utils.GetExternalIP(), hex.EncodeToString(b)), nil
}

func (m *Mutex) getDelay() time.Duration {
	var n int64
	// nolint:errcheck
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	n %= int64(m.delayMax - m.delayMin)
	return time.Duration(n) + m.delayMin
}

func (m *Mutex) acquire(pool Pool, token string) bool {
	client := pool.Get()
	// nolint:errcheck
	ok, err := client.SetNX(m.name, token, m.expiry).Result()
	return err == nil && ok
}

var deleteScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (m *Mutex) release(pool Pool, token string) bool {
	client := pool.Get()
	// nolint:errcheck
	result, err := deleteScript.Run(client, []string{m.name}, token).Int()
	return err == nil && result == 1
}

var touchScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("PEXPIRE", KEYS[1], ARGV[2])
	else
		local rv = redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
		if type(rv) == "table" then
			return 1
		else
			return 0
		end
	end
`)

func (m *Mutex) touch(pool Pool, token string, expiry time.Duration) bool {
	client := pool.Get()
	// nolint:errcheck
	result, err := touchScript.Run(client, []string{m.name}, token, int(expiry/time.Millisecond)).Int()
	return err == nil && result == 1
}
