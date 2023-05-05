/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package abs

import "sync-service/util"

type Sync interface {
	Register(key string, callback func() *util.SafeChannel) string
}
