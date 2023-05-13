/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package ssync

import (
	"errors"
	"fmt"
	"io"

	"sync-service/abs"
	"sync-service/app"
	"sync-service/files/proto/sync_svc"
	"sync-service/util"
)

type SyncService struct {
	sync_svc.UnimplementedSyncServiceServer
}

type srKeyLock struct {
	id          string
	sync        abs.Sync
	isRunning   bool
	errCallback error
	chStatus    *util.SafeChannel
}

func New() *SyncService {
	return &SyncService{}
}

func (*SyncService) notFound(req *sync_svc.KeyLockRequest, id string) *sync_svc.KeyLockResponse {
	return &sync_svc.KeyLockResponse{
		Channel: req.Channel,
		Key:     req.Key,
		Code:    "not-found",
		Message: id,
	}
}

func (*SyncService) execute(req *sync_svc.KeyLockRequest, id string) *sync_svc.KeyLockResponse {
	return &sync_svc.KeyLockResponse{
		Channel: req.Channel,
		Key:     req.Key,
		Code:    "execute",
		Message: id,
	}
}

func (slf *SyncService) KeyLock(stream sync_svc.SyncService_KeyLockServer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover from panic: %+v", r)
		}
	}()

	return slf.keyLock(stream)
}

func (slf *SyncService) keyLock(stream sync_svc.SyncService_KeyLockServer) error {
	defer func() {
		util.Printf("KeyLock: done\n")
	}()

	kl := &srKeyLock{
		isRunning: true,
	}

	for {
		util.Printf("KeyLock, stream.Recv()\n")
		req, err := stream.Recv()
		if kl.errCallback != nil {
			util.Printf("KeyLock [%v], kl.errCallback: %+v\n", kl.id, kl.errCallback)
			kl.stop()
			return kl.errCallback
		}

		if err == io.EOF {
			util.Printf("KeyLock [%v], err io.EOF: %+v\n", kl.id, err)
			kl.stop()
			return nil
		}

		if err != nil {
			util.Printf("KeyLock [%v], err: %+v\n", kl.id, err)
			kl.stop()
			return err
		}

		if kl.sync == nil {
			kl.sync = app.GetChannel(req.Channel)
		}

		util.Printf("KeyLock [%v]: action: %v\n", kl.id, req.Action)

		switch {
		case kl.sync == nil:
			err := stream.Send(slf.notFound(req, kl.id))
			if err != nil {
				kl.stop()
				return err
			}

		case req.Action == "done" && kl.chStatus != nil:
			util.Printf("KeyLock [%v]: kl.chStatus done\n", kl.id)
			kl.stop()

		case req.Action == "register":
			id, executed, panicErr := kl.sync.Register(app.Env.RegisterTimeoutSecond, req.Key, func() *util.SafeChannel {
				defer func() {
					if r := recover(); r != nil {
						util.Printf("recover from panic on callback func, err: %+v", r)
					}
				}()

				if kl.isRunning {
					util.Printf("KeyLock [%v]: send execute\n", kl.id)
					kl.errCallback = stream.Send(slf.execute(req, kl.id))
					util.Printf("KeyLock [%v]: sended execute, isErr: %v", kl.id, kl.errCallback != nil)
					ch := make(chan string)
					chStatus := util.NewSafeChannel(fmt.Sprintf("keylock-register: %v", kl.id), &ch)
					kl.chStatus = chStatus
					return kl.chStatus
				}

				return nil
			})

			if panicErr != nil {
				kl.stop()
				return panicErr
			}

			if !executed {
				kl.stop()
				return errors.New("not executed")
			}

			kl.id = id
			util.Printf("KeyLog register id: %v\n", kl.id)

		default:
			util.Printf("not found switch value\n")
		}
	}
}

func (slf *srKeyLock) stop() {
	slf.isRunning = false
	chStatus := slf.chStatus
	slf.chStatus = nil
	if chStatus != nil {
		chStatus.Send("done")
	}
}
