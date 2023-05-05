/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan. All Rights Reserved.
 *
 * This product is protected by copyright and distributed under
 * licenses restricting copying, distribution and decompilation.
 */

/* cspell: disable-next-line */
package printlog

import (
	"fmt"
	"net/http"
	"strconv"
	"sync-service/app"
	"sync-service/util"

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

	active := v == 1
	app.Env.PrintLog = active
	util.PrintLog = active

	return slf.ctx.RString(nil, http.StatusOK, fmt.Sprintf("success: %v", active))
}
