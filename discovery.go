/**
 * @Author: dingQingHui
 * @Description:
 * @File: func
 * @Version: 1.0.0
 * @Date: 2024/11/18 11:26
 */

package cluster

import (
	"errors"
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"github.com/dingqinghui/extend/component"
	"github.com/dingqinghui/zlog"
	"go.uber.org/zap"
	"time"
)

var (
	errCacheNotInit        = errors.New("cache not init")
	errCacheRequestTimeout = errors.New("request timeout(1s)")
)

type IDiscovery interface {
	component.IComponent
	GetById(kind string) (*common.Member, error)
	GetByKind(kind string) ([]*common.Member, error)
	GetAll() ([]*common.Member, error)
}

type MemberList struct {
	Members []*common.Member
}

func New(options ...Option) IDiscovery {
	d := new(discovery)
	d.opts = loadOptions(options...)
	return d
}

type discovery struct {
	component.BuiltinComponent
	opts *Options
	hub  actor.IProcess
}

func (d *discovery) Name() string {
	return "discovery"
}

func (d *discovery) Init() {
	s := &hubActor{opts: d.opts}
	blueprint := actor.NewBlueprint()
	pid, err := d.opts.system.Spawn(blueprint, func() actor.IActor { return s }, nil)
	if err != nil {
		zlog.Panic("new discovery cache err", zap.Error(err))
		return
	}
	d.hub = pid
}

func (d *discovery) GetById(nodeId string) (*common.Member, error) {
	reply := new(common.Member)
	err := d.requestHub("GetById", nodeId, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (d *discovery) GetByKind(kind string) ([]*common.Member, error) {
	reply := new(MemberList)
	err := d.requestHub("GetByKind", kind, reply)
	if err != nil {
		return nil, err
	}
	return reply.Members, err
}

func (d *discovery) GetAll() ([]*common.Member, error) {
	reply := new(MemberList)
	err := d.requestHub("GetAll", nil, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Members, err
}

func (d *discovery) requestHub(funcName string, arg, reply interface{}) error {
	if d.hub == nil {
		zlog.Error("discovery call cache actor err", zap.String("funcName", funcName), zap.Error(errCacheNotInit))
		return errCacheNotInit
	}
	err := d.hub.Call(funcName, time.Second, arg, reply)
	if err != nil {
		zlog.Error("discovery call cache actor err", zap.String("funcName", funcName), zap.Error(err))
		return err
	}
	return nil
}

func (d *discovery) Stop() {
	if d.hub == nil {
		return
	}
	_ = d.hub.Stop()
}
