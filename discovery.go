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
	GetById(kind string) (common.IMember, error)
	GetByKind(kind string) ([]common.IMember, error)
	GetAll() ([]common.IMember, error)
}

func New(options ...Option) IDiscovery {
	d := new(discovery)
	d.opts = loadOptions(options...)
	return d
}

type discovery struct {
	component.BuiltinComponent
	opts     *Options
	cachePid actor.IProcess
}

func (d *discovery) Name() string {
	return "discovery"
}

func (d *discovery) Init() {
	s := &cacheActor{opts: d.opts}
	blueprint := actor.NewBlueprint()
	pid, err := d.opts.system.Spawn(blueprint, func() actor.IActor { return s }, nil)
	if err != nil {
		zlog.Panic("new discovery cache err", zap.Error(err))
		return
	}
	d.cachePid = pid
}

func (d *discovery) Stop() {
	if d.cachePid == nil {
		return
	}
	_ = d.cachePid.Stop()
}

func (d *discovery) GetById(kind string) (common.IMember, error) {
	res, err := d.requestCache("GetKind", kind)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(common.IMember), nil
}

func (d *discovery) GetByKind(kind string) ([]common.IMember, error) {
	res, err := d.requestCache("GetByKind", kind)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.([]common.IMember), nil
}

func (d *discovery) GetAll() ([]common.IMember, error) {
	res, err := d.requestCache("GetAll", "")
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.([]common.IMember), nil
}

func (d *discovery) requestCache(funcName string, message interface{}) (interface{}, error) {
	if d.cachePid == nil {
		zlog.Error("discovery call cache actor err", zap.String("funcName", funcName), zap.Error(errCacheNotInit))
		return nil, errCacheNotInit
	}
	res, isTimeout, err := d.cachePid.Call(funcName, message, time.Second)
	if err != nil {
		return nil, err
	}
	if isTimeout {
		zlog.Error("discovery call cache actor err", zap.String("funcName", funcName), zap.Error(errCacheRequestTimeout))
		return nil, errCacheRequestTimeout
	}
	return res, nil
}
