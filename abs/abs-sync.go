/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package abs

import (
	"sync-service/util"
	"time"
)

type Sync interface {
	Register(timeout time.Duration, key string, callback func() *util.SafeChannel) (id string, executed bool, panicErr error)
}
