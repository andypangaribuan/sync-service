/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package app

import (
	"fmt"
	"sync"
	"sync-service/util"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

type srSync struct {
	channel string
	mutex   sync.Mutex
	pools   map[string]*srPool
}

type srPool struct {
	key       string
	mtx       *sync.RWMutex
	isActive  bool
	chStop    chan bool
	callbacks map[string]func() *util.SafeChannel
}

func (slf *srSync) generateId() string {
	return fmt.Sprintf("%v%v", f9.TimeNow().UnixMilli(), p9.Util.GetNanoID(10))
}

func (slf *srSync) cleansing(exceptKey string) {
	length := len(slf.pools)
	if length > 0 {
		rm := make([]string, 0)
		for k, v := range slf.pools {
			if len(v.callbacks) == 0 && v.key != exceptKey {
				rm = append(rm, k)
			}
		}

		for _, key := range rm {
			pool := slf.pools[key]
			pool.isActive = false
			<-pool.chStop
			close(pool.chStop)
			delete(slf.pools, key)
		}
	}
}

func (slf *srSync) Register(key string, callback func() *util.SafeChannel) string {
	slf.mutex.Lock()
	defer slf.mutex.Unlock()

	util.Printf("register start, key: %v\n", key)
	defer util.Printf("register end, key: %v\n", key)

	slf.cleansing(key)
	iteration := 0

	for {
		iteration++
		length := len(slf.pools)

		if length >= Env.MutexSizePerChannel {
			time.Sleep(time.Millisecond * 10)
		} else if iteration >= 30 {
			iteration = 0
			slf.cleansing(key)
		} else {
			break
		}
	}

	id := slf.generateId()
	util.Printf("register id: %v\n", id)
	pool, ok := slf.pools[key]

	if !ok {
		pool = &srPool{
			key: key,
			mtx: &sync.RWMutex{},
			callbacks: map[string]func() *util.SafeChannel{
				id: callback,
			},
			isActive: true,
			chStop:   make(chan bool),
		}

		slf.pools[key] = pool
		util.Printf("register run: %v\n", id)
		go pool.runner()
	} else {
		pool.mtx.Lock()
		defer pool.mtx.Unlock()

		pool.callbacks[id] = callback
	}

	return id
}

func (slf *srPool) runner() {
	iteration := 0

	for {
		if !slf.isActive {
			break
		}

		time.Sleep(time.Millisecond * 100)
		iteration++

		if iteration >= Env.RunnerIteration {
			iteration = 0
			util.Printf("runner: %v [active]\n", slf.key)
		}

		slf.run()
	}

	slf.chStop <- true
}

func (slf *srPool) run() {
	ls := make(map[string]func() *util.SafeChannel, 0)

	slf.mtx.RLock()
	for k, v := range slf.callbacks {
		ls[k] = v
	}
	slf.mtx.RUnlock()

	ids := make([]string, 0)

	for id, callback := range ls {
		util.Printf("- run(): call callback [%v]\n", id)
		chStatus := callback()
		if chStatus == nil {
			ids = append(ids, id)
		} else {
			util.Printf("- run(): wait chStatus [%v]\n", id)
			status := chStatus.Read("done")
			util.Printf("- run(): status %v [%v]\n", status, id)
			if status == "done" {
				ids = append(ids, id)
				chStatus.Close()
				util.Printf("- run(): closed chStatus [%v]\n", id)
			}
		}
	}

	if len(ids) > 0 {
		slf.mtx.Lock()
		defer slf.mtx.Unlock()

		for _, id := range ids {
			delete(slf.callbacks, id)
		}
	}

}

func (slf *srSync) check() {
	for {
		time.Sleep(Env.PoolCheckDelay)

		totalPool := len(slf.pools)
		msg := ""

		i := 0
		for k, v := range slf.pools {
			i++
			if msg != "" {
				msg += "\n"
			}

			if v == nil {
				msg += fmt.Sprintf("%v. %v: [NIL]", i, k)
			} else {
				msg += fmt.Sprintf("%v. %v: %v", i, k, len(v.callbacks))
			}
		}

		if msg == "" {
			util.Printf("total pool: %v\n", totalPool)
		} else {
			util.Printf("total pool: %v\n%v\n", totalPool, msg)
		}
	}
}
