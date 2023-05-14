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
	"sync-service/abs"

	"github.com/andypangaribuan/project9"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
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
		sc := &srSync{
			channel: channel,
			mtx:     p9.Util.NewMutex(fmt.Sprintf("channel: %v", channel)),
			pools:   make(map[string]*srPool, 0),
		}

		syncChannels[channel] = sc
		go sc.check()
		go sc.autoClean()
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
