package ipblocker

import (
	"sync"
	"time"
)

// a ConcurrentDictionary
var rate *rateLimit

func init() {
	rate = newRateLimit()
}

// Success set the ips struct to refresh
func Success(ip string) {
	rate.Lock()
	defer rate.Unlock()
	result := rate.GetorAdd(ip, newLockInfo())
	result.Success()
}

// add one try count to this ips struct
func Fail(ip string) {
	rate.Lock()
	defer rate.Unlock()
	result := rate.GetorAdd(ip, newLockInfo())
	result.Fail()
}

// check the ips status and return it can log in  or not
func IsLoginable(ip string) bool {
	rate.Lock()
	defer rate.Unlock()
	result := rate.GetorAdd(ip, newLockInfo())
	return result.IsLoginable()
}

// the interface for real ips struct
type limits interface {
	Fail()
	Success()
	IsLoginable() bool
}

// a ConcurrentDictionary
type rateLimit struct {
	RateMap map[string]limits
	sync.Mutex
}

func newRateLimit() *rateLimit {
	return &rateLimit{
		RateMap: make(map[string]limits),
	}
}

// a useful function in c#, if it doesn't have one, then add one
func (rate *rateLimit) GetorAdd(ip string, limit limits) limits {
	if l, ok := rate.RateMap[ip]; ok {
		return l
	} else {
		rate.RateMap[ip] = limit
		return limit
	}
}

// the struct for block ip
type lockInfo struct {
	TryCount    int
	BlockedTime time.Time
}

func newLockInfo() *lockInfo {
	return &lockInfo{
		TryCount:    0,
		BlockedTime: time.Time{},
	}
}
func (lock *lockInfo) Fail() {
	lock.TryCount++
	if lock.TryCount > 5 {
		lock.TryCount = 0
		lock.BlockedTime = time.Now()
	}
}

func (lock *lockInfo) Success() {
	lock.TryCount = 0
}

func (lock *lockInfo) IsLoginable() bool {
	return lock.BlockedTime.Add(5*time.Minute).Before(time.Now()) && lock.TryCount <= 5
}
