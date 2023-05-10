package registry

import (
	"context"
	"kit.golaxy.org/golaxy/service"
	"time"
)

// Register 注册服务
func Register(serviceCtx service.Context, ctx context.Context, service Service, ttl time.Duration) error {
	return Get(serviceCtx).Register(ctx, service, ttl)
}

// Deregister 取消注册服务
func Deregister(serviceCtx service.Context, ctx context.Context, service Service) error {
	return Get(serviceCtx).Deregister(ctx, service)
}

// GetService 查询服务
func GetService(serviceCtx service.Context, ctx context.Context, serviceName string) ([]Service, error) {
	return Get(serviceCtx).GetService(ctx, serviceName)
}

// ListServices 查询所有服务
func ListServices(serviceCtx service.Context, ctx context.Context) ([]Service, error) {
	return Get(serviceCtx).ListServices(ctx)
}

// Watch 获取服务监听器
func Watch(serviceCtx service.Context, ctx context.Context, serviceName string) (Watcher, error) {
	return Get(serviceCtx).Watch(ctx, serviceName)
}
