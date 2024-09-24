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

package etcd_dsync

import (
	"context"
	"crypto/tls"
	"fmt"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/option"
	"git.golaxy.org/framework/net/netpath"
	"git.golaxy.org/framework/plugins/dsync"
	"git.golaxy.org/framework/plugins/log"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func newDSync(settings ...option.Setting[DSyncOptions]) dsync.IDistSync {
	return &_DistSync{
		options: option.Make(With.Default(), settings...),
	}
}

type _DistSync struct {
	svcCtx  service.Context
	options DSyncOptions
	client  *etcdv3.Client
}

// InitSP 初始化服务插件
func (s *_DistSync) InitSP(svcCtx service.Context) {
	log.Infof(svcCtx, "init plugin %q", self.Name)

	s.svcCtx = svcCtx

	if s.options.EtcdClient == nil {
		cli, err := etcdv3.New(s.configure())
		if err != nil {
			log.Panicf(svcCtx, "new etcd client failed, %s", err)
		}
		s.client = cli
	} else {
		s.client = s.options.EtcdClient
	}

	for _, ep := range s.client.Endpoints() {
		func() {
			ctx, cancel := context.WithTimeout(s.svcCtx, 3*time.Second)
			defer cancel()

			if _, err := s.client.Status(ctx, ep); err != nil {
				log.Panicf(s.svcCtx, "status etcd %q failed, %s", ep, err)
			}
		}()
	}
}

// ShutSP 关闭服务插件
func (s *_DistSync) ShutSP(svcCtx service.Context) {
	log.Infof(svcCtx, "shut plugin %q", self.Name)

	if s.options.EtcdClient == nil {
		if s.client != nil {
			s.client.Close()
		}
	}
}

// NewMutex returns a new distributed mutex with given name.
func (s *_DistSync) NewMutex(name string, settings ...option.Setting[dsync.DistMutexOptions]) dsync.IDistMutex {
	return s.newMutex(name, option.Make(dsync.With.Default(), settings...))
}

// NewMutexf returns a new distributed mutex using a formatted string.
func (s *_DistSync) NewMutexf(format string, args ...any) dsync.IDistMutexSettings {
	return &_DistMutexSettings{
		dsync: s,
		name:  fmt.Sprintf(format, args...),
	}
}

// NewMutexp returns a new distributed mutex using elements.
func (s *_DistSync) NewMutexp(elems ...string) dsync.IDistMutexSettings {
	return &_DistMutexSettings{
		dsync: s,
		name:  netpath.Join(s.GetSeparator(), elems...),
	}
}

// GetSeparator return name path separator.
func (s *_DistSync) GetSeparator() string {
	return "/"
}

func (s *_DistSync) configure() etcdv3.Config {
	if s.options.EtcdConfig != nil {
		return *s.options.EtcdConfig
	}

	config := etcdv3.Config{
		Endpoints:   s.options.CustomAddresses,
		Username:    s.options.CustomUsername,
		Password:    s.options.CustomPassword,
		DialTimeout: 3 * time.Second,
	}

	if s.options.CustomTLSConfig != nil {
		tlsConfig := s.options.CustomTLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig
	}

	return config
}
