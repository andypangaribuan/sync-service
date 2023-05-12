/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package app

import (
	"strings"
	"time"

	"github.com/andypangaribuan/project9/act/actenv"
	"github.com/andypangaribuan/project9/p9"
)

type srEnv struct {
	AppName        string
	AppVersion     string
	AppPortRestful int
	AppPortGRPC    int
	AppEnv         actenv.AppEnv
	AppAutoRecover bool
	Timezone       string

	Channels            []string
	MutexSizePerChannel int

	PrintLog                    bool
	RunnerIteration             int
	PoolCheckDelay              time.Duration
	RegisterTimeoutSecond       time.Duration
	StreamCallbackTimeoutSecond time.Duration
}

func envInitialize() {
	poolCheckDelay := p9.Util.Env.GetInt("POOL_CHECK_DELAY", 10)
	if poolCheckDelay <= 0 {
		poolCheckDelay = 10
	}

	Env = &srEnv{
		AppName:        p9.Util.Env.GetStr("APP_NAME"),
		AppVersion:     p9.Util.Env.GetStr("APP_VERSION", "0.0.0"),
		AppPortRestful: p9.Util.Env.GetInt("APP_PORT_RESTFUL"),
		AppPortGRPC:    p9.Util.Env.GetInt("APP_PORT_GRPC"),
		AppEnv:         p9.Util.Env.GetAppEnv("APP_ENV"),
		AppAutoRecover: p9.Util.Env.GetBool("APP_AUTO_RECOVER"),
		Timezone:       p9.Util.Env.GetStr("TIMEZONE"),

		Channels:            strings.Split(p9.Util.Env.GetStr("CHANNELS"), ","),
		MutexSizePerChannel: p9.Util.Env.GetInt("MUTEX_SIZE_PER_CHANNEL"),

		PrintLog:                    p9.Util.Env.GetBool("PRINT_LOG", false),
		RunnerIteration:             p9.Util.Env.GetInt("RUNNER_ITERATION", 100),
		PoolCheckDelay:              time.Second * time.Duration(poolCheckDelay),
		RegisterTimeoutSecond:       time.Second * time.Duration(p9.Util.Env.GetInt64("REGISTER_TIMEOUT_SECOND")),
		StreamCallbackTimeoutSecond: time.Second * time.Duration(p9.Util.Env.GetInt64("STREAM_CALLBACK_TIMEOUT_SECOND")),
	}
}
