/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

/* cspell: disable-next-line */
package runneriteration

import (
	"net/http"
	"strconv"
	"sync-service/app"

	"github.com/andypangaribuan/project9/server"
)

type srController struct {
	ctx server.FuseContext
}

func New(ctx server.FuseContext) func() error {
	sr := &srController{
		ctx: ctx,
	}
	return sr.handler
}

func (slf *srController) handler() error {
	return slf.process()
}

func (slf *srController) process() error {
	val := slf.ctx.Query("val", "")
	if val == "" {
		return slf.ctx.RString(nil, http.StatusBadRequest, "required: query val [int]")
	}

	v, err := strconv.Atoi(val)
	if err != nil {
		return slf.ctx.RString(nil, http.StatusBadRequest, "invalid: query val [int]")
	}

	if v <= 0 {
		return slf.ctx.RString(nil, http.StatusBadRequest, "invalid: query val [int] > 0")
	}

	app.Env.RunnerIteration = v
	return slf.ctx.RString(nil, http.StatusOK, "success")
}
