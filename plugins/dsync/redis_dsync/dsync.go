package redis_dsync

import (
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/option"
	"git.golaxy.org/framework/plugins/dsync"
	"git.golaxy.org/framework/plugins/log"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func newDSync(settings ...option.Setting[DSyncOptions]) dsync.IDistSync {
	return &_DistSync{
		options: option.Make(With.Default(), settings...),
	}
}

type _DistSync struct {
	options DSyncOptions
	servCtx service.Context
	client  *redis.Client
	redSync *redsync.Redsync
}

// InitSP 初始化服务插件
func (s *_DistSync) InitSP(ctx service.Context) {
	log.Infof(ctx, "init plugin %q", self.Name)

	s.servCtx = ctx

	if s.options.RedisClient == nil {
		s.client = redis.NewClient(s.configure())
	} else {
		s.client = s.options.RedisClient
	}

	_, err := s.client.Ping(ctx).Result()
	if err != nil {
		log.Panicf(ctx, "ping redis %q failed, %v", s.client, err)
	}

	s.redSync = redsync.New(goredis.NewPool(s.client))
}

// ShutSP 关闭服务插件
func (s *_DistSync) ShutSP(ctx service.Context) {
	log.Infof(ctx, "shut plugin %q", self.Name)

	if s.options.RedisClient == nil {
		if s.client != nil {
			s.client.Close()
		}
	}
}

// NewMutex returns a new distributed mutex with given name.
func (s *_DistSync) NewMutex(name string, settings ...option.Setting[dsync.DistMutexOptions]) dsync.IDistMutex {
	return s.newMutex(name, option.Make(dsync.With.Default(), settings...))
}

// GetSeparator return name path separator.
func (s *_DistSync) GetSeparator() string {
	return ":"
}

func (s *_DistSync) configure() *redis.Options {
	if s.options.RedisConfig != nil {
		return s.options.RedisConfig
	}

	if s.options.RedisURL != "" {
		conf, err := redis.ParseURL(s.options.RedisURL)
		if err != nil {
			log.Panicf(s.servCtx, "parse redis url %q failed, %s", s.options.RedisURL, err)
		}
		return conf
	}

	conf := &redis.Options{}
	conf.Username = s.options.CustomUsername
	conf.Password = s.options.CustomPassword
	conf.Addr = s.options.CustomAddress
	conf.DB = s.options.CustomDB

	return conf
}
