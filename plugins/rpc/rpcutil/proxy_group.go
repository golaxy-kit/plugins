/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package rpcutil

import (
	"errors"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/framework/plugins/gate"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/callpath"
	"git.golaxy.org/framework/plugins/rpcstack"
	"strings"
)

// ProxyGroup 代理分组
func ProxyGroup(ctx runtime.CurrentContextProvider, name string) GroupProxied {
	return GroupProxied{
		servCtx: service.Current(ctx),
		rtCtx:   runtime.Current(ctx),
		addr:    gate.CliDetails.MulticastSubdomainJoin(name),
	}
}

// ConcurrentProxyGroup 代理分组
func ConcurrentProxyGroup(ctx service.Context, name string) GroupProxied {
	return GroupProxied{
		servCtx: ctx,
		addr:    gate.CliDetails.MulticastSubdomainJoin(name),
	}
}

// GroupProxied 分组代理，用于向分组发送RPC
type GroupProxied struct {
	servCtx service.Context
	rtCtx   runtime.Context
	addr    string
}

// GetName 获取分组名称
func (p GroupProxied) GetName() string {
	return strings.TrimPrefix(strings.TrimPrefix(p.addr, gate.CliDetails.MulticastSubdomain), gate.CliDetails.PathSeparator)
}

// GetAddr 获取分组地址
func (p GroupProxied) GetAddr() string {
	return p.addr
}

// OneWayCliRPC 向分组发送单向RPC
func (p GroupProxied) OneWayCliRPC(method string, args ...any) error {
	return p.OneWayCliRPCToEntity(uid.Nil, method, args...)
}

// OneWayCliRPCToEntity 向分组发送单向RPC
func (p GroupProxied) OneWayCliRPCToEntity(entityId uid.Id, method string, args ...any) error {
	if p.servCtx == nil {
		panic(errors.New("rpc: setting servCtx is nil"))
	}

	// 调用链
	callChain := rpcstack.EmptyCallChain
	if p.rtCtx != nil {
		callChain = rpcstack.Using(p.rtCtx).CallChain()
	}

	// 调用路径
	cp := callpath.CallPath{
		Category: callpath.Client,
		EntityId: entityId,
		Method:   method,
	}

	return rpc.Using(p.servCtx).OneWayRPC(p.addr, callChain, cp.String(), args...)
}
