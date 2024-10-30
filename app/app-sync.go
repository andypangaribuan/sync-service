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
	"log"
	"sync-service/util"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

type srSync struct {
	channel string
	mtx     abs.UtilMutex
	pools   map[string]*srPool
}

type srPool struct {
	key       string
	mtx       abs.UtilMutex
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

func (slf *srSync) Register(timeout time.Duration, key string, callback func() *util.SafeChannel) (id string, executed bool, panicErr error) {
	// slf.mtx.Lock()
	// defer slf.mtx.Unlock()

	executed, panicErr = slf.mtx.Exec(&timeout, func() {
		id = slf.register(key, callback)
	})

	return
}

func (slf *srSync) register(key string, callback func() *util.SafeChannel) string {
	util.Printf("register start, key: %v\n", key)
	defer util.Printf("register end, key: %v\n", key)

	// log.Printf("do cleansing, key: %v\n", key)
	slf.cleansing(key)
	iteration := 0
	const maxIteration = 10

	// log.Printf("do iteration, key: %v\n", key)
	for {
		iteration++
		length := len(slf.pools)

		switch {
		case length >= Env.MutexSizePerChannel && iteration < maxIteration:
			time.Sleep(time.Millisecond * 10)

		case iteration >= maxIteration:
			iteration = 0
			slf.cleansing(key)
		}

		if length < Env.MutexSizePerChannel {
			break
		}

		// if length >= Env.MutexSizePerChannel {
		// 	time.Sleep(time.Millisecond * 10)
		// } else if iteration >= 30 {
		// 	iteration = 0
		// 	slf.cleansing(key)
		// } else {
		// 	break
		// }
	}

	id := slf.generateId()
	util.Printf("register id: %v\n", id)
	pool, ok := slf.pools[key]

	if !ok {
		pool = &srPool{
			key: key,
			mtx: p9.Util.NewMutex(fmt.Sprintf("pool: %v", key)),
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
		log.Printf("add to pool, key: %v, id: %v\n", key, id)
		pool.mtx.Exec(nil, func() {
			pool.callbacks[id] = callback
		})
		log.Printf("end to pool, key: %v, id: %v\n", key, id)
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

	slf.mtx.Exec(nil, func() {
		for k, v := range slf.callbacks {
			ls[k] = v
		}
	})

	ids := make([]string, 0)

	for id, callback := range ls {
		util.Printf("- run(): call callback [%v]\n", id)
		chStatus := callback()

		if chStatus == nil {
			ids = append(ids, id)
		} else {
			util.Printf("- run(): wait chStatus [%v]\n", id)
			status := chStatus.Read(Env.StreamCallbackTimeoutSecond, "done")
			util.Printf("- run(): status %v [%v]\n", status, id)

			if status == "done" {
				ids = append(ids, id)
				chStatus.Close()
				util.Printf("- run(): closed chStatus [%v]\n", id)
			}
		}
	}

	if len(ids) > 0 {
		slf.mtx.Exec(nil, func() {
			for _, id := range ids {
				delete(slf.callbacks, id)
			}
		})
	}
}

func (slf *srSync) autoClean() {
	for {
		time.Sleep(Env.AutoCleanPoolDelaySecond)

		slf.mtx.Exec(&Env.AutoCleanPoolDelaySecond, func() {
			slf.cleansing("")
		})
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
				msg += fmt.Sprintf("%v. [%v] %v: [NIL]", i, slf.channel, k)
			} else {
				msg += fmt.Sprintf("%v. [%v] %v: %v", i, slf.channel, k, len(v.callbacks))
			}
		}

		if msg == "" {
			util.Printf("total pool [%v]: %v\n", slf.channel, totalPool)
		} else {
			util.Printf("total pool [%v]: %v\n%v\n", slf.channel, totalPool, msg)
		}
	}
}
