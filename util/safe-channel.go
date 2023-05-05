/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package util

import (
	"sync"

	"github.com/andypangaribuan/project9/f9"
)

type SafeChannel struct {
	CH     *chan string
	closed bool
	mtx    sync.Mutex
}

func (slf *SafeChannel) Send(val string) {
	if slf.closed {
		return
	}

	slf.mtx.TryLock()
	defer slf.mtx.Unlock()

	if !slf.closed && slf.CH != nil {
		*slf.CH <- val
	}
}

func (slf *SafeChannel) Read(defaultValue ...string) string {
	if !slf.closed && slf.CH != nil {
		return <-*slf.CH
	}

	return f9.Ternary(len(defaultValue) == 0, "", defaultValue[0])
}

func (slf *SafeChannel) Close() {
	if slf.closed {
		return
	}

	slf.mtx.TryLock()
	defer slf.mtx.Unlock()

	if !slf.closed && slf.CH != nil {
		slf.closed = true
		close(*slf.CH)
	}
}
