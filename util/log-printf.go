/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package util

import (
	"log"
)

var PrintLog bool

func init() {
	PrintLog = false
}

func Printf(format string, arg ...any) {
	if PrintLog {
		log.Printf(format, arg...)
	}
}
