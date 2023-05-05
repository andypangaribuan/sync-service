/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package main

import (
	"log"
	"sync-service/app"
	"sync-service/ctrl"
	"sync-service/files/proto/sync_svc"
	"sync-service/util"

	"github.com/andypangaribuan/project9/server"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	withErr := !app.Env.AppEnv.IsProd() && !app.Env.AppEnv.IsSandbox()
	util.PrintLog = app.Env.PrintLog
	server.Fuse2(app.Env.AppPortRestful, app.Env.AppPortGRPC, app.Env.AppAutoRecover, withErr, routes, register)
}

func routes(router server.FuseRouter) {
	router.Group(map[string][]func(sc server.FuseContext) error{
		"GET: /private/svc-name":    {ctrl.Private.GetSvcName},
		"GET: /private/svc-version": {ctrl.Private.GetSvcVersion},
		"GET: /private/svc-check":   {ctrl.Private.SvcCheck},

		"GET: /private/pool-check":       {ctrl.Private.PoolCheck},
		"GET: /private/runner-iteration": {ctrl.Private.RunnerIteration},
		"GET: /private/print-log":        {ctrl.Private.PrintLog},
	})
}

func register(server *grpc.Server) {
	sync_svc.RegisterSyncServiceServer(server, ctrl.SyncSvc)
}
