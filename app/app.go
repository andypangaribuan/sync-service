/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package app

import (
	"sync"
	"sync-service/abs"

	"github.com/andypangaribuan/project9"
	"github.com/andypangaribuan/project9/f9"
)

var (
	Env          *srEnv
	syncChannels map[string]*srSync
)

func init() {
	project9.Initialize()

	envInitialize()

	syncChannels = make(map[string]*srSync, 0)
	for i := 0; i < len(Env.Channels); i++ {
		channel := Env.Channels[i]
		syncChannels[channel] = &srSync{
			channel: channel,
			mutex:   sync.Mutex{},
			pools:   make(map[string]*srPool, 0),
		}

		go syncChannels[channel].check()
	}

	f9.TimeZone = Env.Timezone
}

func GetChannel(channel string) abs.Sync {
	s, ok := syncChannels[channel]
	if ok {
		return s
	}

	return nil
}