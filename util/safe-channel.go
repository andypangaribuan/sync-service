/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package util

import (
	"context"
	"fmt"
	"time"

	"github.com/andypangaribuan/project9/abs"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
)

type SafeChannel struct {
	ch     *chan string
	closed bool
	mtx    abs.UtilMutex
}

func NewSafeChannel(name string, ch *chan string) *SafeChannel {
	return &SafeChannel{
		ch:  ch,
		mtx: p9.Util.NewMutex(name),
	}
}

func (slf *SafeChannel) Send(val string) {
	if slf.closed || slf.ch == nil {
		return
	}

	for {
		executed, panicErr := slf.mtx.Exec(nil, func() {
			if !slf.closed && slf.ch != nil {
				*slf.ch <- val
			}
		})

		if executed || panicErr != nil {
			break
		}
	}

	// slf.mtx.Lock()
	// defer slf.mtx.Unlock()

	// if !slf.closed && slf.ch != nil {
	// 	*slf.ch <- val
	// }
}

func (slf *SafeChannel) Read(timeout time.Duration, defaultValue ...string) (val string) {
	if !slf.closed && slf.ch != nil {
		var panicErr error
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		go func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					panicErr = fmt.Errorf("panic error: %+v", r)
				}
				cancel()
			}()

			val = <-*slf.ch
		}(ctx)

		<-ctx.Done()
		if panicErr != nil {
			return f9.Ternary(len(defaultValue) == 0, "", defaultValue[0])
		}

		// switch ctx.Err() {
		// case context.DeadlineExceeded:
		// 	isTimeout = true
		// }

		// val = <-*slf.ch
		return
	}

	return f9.Ternary(len(defaultValue) == 0, "", defaultValue[0])
}

func (slf *SafeChannel) Close() {
	if slf.closed || slf.ch == nil {
		return
	}

	panicErrCount := 0

	for {
		executed, panicErr := slf.mtx.Exec(nil, func() {
			slf.closed = true
			close(*slf.ch)
		})

		if executed {
			break
		}

		if panicErr != nil {
			panicErrCount++
		}

		if panicErrCount >= 3 {
			break
		}
	}

	// slf.mtx.Lock()
	// defer slf.mtx.Unlock()

	// if !slf.closed && slf.ch != nil {
	// 	slf.closed = true
	// 	close(*slf.ch)
	// }
}
