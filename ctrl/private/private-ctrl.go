/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

/* spell-checker: disable */
package private

import (
	getsvcname "sync-service/ctrl/private/get-svc-name"
	getsvcversion "sync-service/ctrl/private/get-svc-version"
	poolcheck "sync-service/ctrl/private/pool-check"
	printlog "sync-service/ctrl/private/print-log"
	runneriteration "sync-service/ctrl/private/runner-iteration"
	svccheck "sync-service/ctrl/private/svc-check"

	"github.com/andypangaribuan/project9/server"
)

type PrivateCtrl struct{}

func (slf *PrivateCtrl) GetSvcName(ctx server.FuseContext) error {
	return getsvcname.New(ctx)()
}

func (slf *PrivateCtrl) GetSvcVersion(ctx server.FuseContext) error {
	return getsvcversion.New(ctx)()
}

func (slf *PrivateCtrl) SvcCheck(ctx server.FuseContext) error {
	return svccheck.New(ctx)()
}

func (slf *PrivateCtrl) PoolCheck(ctx server.FuseContext) error {
	return poolcheck.New(ctx)()
}

func (slf *PrivateCtrl) RunnerIteration(ctx server.FuseContext) error {
	return runneriteration.New(ctx)()
}

func (slf *PrivateCtrl) PrintLog(ctx server.FuseContext) error {
	return printlog.New(ctx)()
}
