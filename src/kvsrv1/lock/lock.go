package lock

import (
	"time"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck kvtest.IKVClerk

	lockKey   string
	lockValue string
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{
		ck:        ck,
		lockKey:   l,
		lockValue: kvtest.RandValue(8),
	}
	return lk
}

func (lk *Lock) Acquire() {

	for {
		val, ver, err := lk.ck.Get(lk.lockKey)

		if err == rpc.ErrNoKey {
			err = lk.ck.Put(lk.lockKey, lk.lockValue, 0)

			if err == rpc.OK {
				return
			}
		}

		if err != rpc.OK {
			continue
		}

		if val == "" {
			err = lk.ck.Put(lk.lockKey, lk.lockValue, ver)
			if err == rpc.OK {
				return
			}
		} else if val == lk.lockValue {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (lk *Lock) Release() {
	for {
		val, ver, err := lk.ck.Get(lk.lockKey)

		if err != rpc.OK {
			continue
		}

		if val == lk.lockValue {
			err = lk.ck.Put(lk.lockKey, "", ver)
			if err == rpc.OK {
				return
			}
		} else if val == "" {
			return
		}

		time.Sleep(100 * time.Millisecond)
	}
}
