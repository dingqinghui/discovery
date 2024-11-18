/**
 * @Author: dingQingHui
 * @Description:
 * @File: options
 * @Version: 1.0.0
 * @Date: 2024/11/18 10:05
 */

package cluster

import (
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"github.com/dingqinghui/discovery/provider/consul"
	"time"
)

type Option func(*Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	opts.provider, _ = consul.NewConsulProvider("127.0.0.1:8500")
	opts.system = actor.NewSystem()
	opts.service = &common.BuiltinMember{
		Name:    "default",
		Id:      "Id",
		Address: "127.0.0.1:8081",
		Port:    3000,
		Tags:    []string{"Tags"},
		Meta:    nil,
	}
	opts.ttl = time.Second * 3
	opts.deregisterTtl = time.Second * 30
	for _, option := range options {
		option(opts)
	}
	return opts
}

type Options struct {
	provider      common.IProvider
	knowKinds     []string
	system        actor.ISystem
	service       common.IMember
	ttl           time.Duration
	deregisterTtl time.Duration
}

func WithHealthDeregisterTtl(deregisterTtl time.Duration) Option {
	return func(op *Options) {
		op.deregisterTtl = deregisterTtl
	}
}
func WithHealthTtl(ttl time.Duration) Option {
	return func(op *Options) {
		op.ttl = ttl
	}
}
func WithRegisterService(service common.IMember) Option {
	return func(op *Options) {
		op.service = service
	}
}

func WithKnowKinds(knowKinds []string) Option {
	return func(op *Options) {
		op.knowKinds = knowKinds
	}
}

func WithProvider(provider common.IProvider) Option {
	return func(op *Options) {
		op.provider = provider
	}
}
func WithSystem(system actor.ISystem) Option {
	return func(op *Options) {
		op.system = system
	}
}
