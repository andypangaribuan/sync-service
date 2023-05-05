/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

package ssync

import (
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

func (*SyncService) waiting(req *sync_svc.KeyLockRequest, id string) *sync_svc.KeyLockResponse {
	return &sync_svc.KeyLockResponse{
		Channel: req.Channel,
		Key:     req.Key,
		Code:    "waiting",
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

func (slf *SyncService) KeyLock(stream sync_svc.SyncService_KeyLockServer) error {
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
			kl.done()

		case req.Action == "register":
			kl.id = kl.sync.Register(req.Key, func() *util.SafeChannel {
				if kl.isRunning {
					util.Printf("KeyLock [%v]: send execute\n", kl.id)
					kl.errCallback = stream.Send(slf.execute(req, kl.id))
					util.Printf("KeyLock [%v]: sended execute, isErr: %v", kl.id, kl.errCallback != nil)
					ch := make(chan string)
					chStatus := &util.SafeChannel{
						CH: &ch,
					}
					kl.chStatus = chStatus
					return kl.chStatus
				}

				return nil
			})
			util.Printf("KeyLog register id: %v\n", kl.id)

		default:
			util.Printf("not found switch value\n")
		}
	}
}

func (slf *srKeyLock) stop() {
	slf.isRunning = false
	slf.done()
}

func (slf *srKeyLock) done() {
	chStatus := slf.chStatus
	slf.chStatus = nil
	if chStatus != nil {
		chStatus.Send("done")
	}
}
