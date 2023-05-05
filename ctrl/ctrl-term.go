/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package ctrl

import (
	"sync-service/ctrl/private"
	"sync-service/ctrl/ssync"
	"sync-service/files/proto/sync_svc"
)

var (
	SyncSvc sync_svc.SyncServiceServer
	Private *private.PrivateCtrl
)

func init() {
	SyncSvc = ssync.New()
	Private = new(private.PrivateCtrl)
}
