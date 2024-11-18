/**
 * @Author: dingQingHui
 * @Description:
 * @File: consul_provider
 * @Version: 1.0.0
 * @Date: 2024/4/25 14:09
 */

package consul

import (
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"github.com/dingqinghui/extend/component"
	"github.com/dingqinghui/zlog"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"time"
)

type consulProvider struct {
	component.BuiltinComponent
	client       *api.Client
	kvWaitIndex  uint64
	healthPid    actor.IProcess
	discoveryPid actor.IProcess
}

func NewConsulProvider(consulAddress string) (common.IProvider, error) {
	c := new(consulProvider)
	if err := c.connect(consulAddress); err != nil {
		zlog.Error("consul connect err", zap.Error(err))
		return nil, err
	}
	return c, nil
}
func (c *consulProvider) connect(consulAddress string) error {
	apiConfig := api.DefaultConfig()
	apiConfig.Address = consulAddress
	client, err := api.NewClient(apiConfig)
	if err != nil {
		zlog.Error("consul new client err", zap.Error(err))
		return err
	}
	c.client = client
	zlog.Info("consul connect success..", zap.String("consulAddress", consulAddress))
	return nil
}

func (c *consulProvider) Name() string {
	return "consul"
}

func (c *consulProvider) Discovery(system actor.ISystem, clusterName string, knowKinds []string, f common.EventMemberUpdateHandler) error {
	if c.discoveryPid != nil {
		return nil
	}
	if len(knowKinds) <= 0 {
		zlog.Warn("consul discovery kinds is nil")
		return nil
	}
	d := &discoveryActor{client: c.client, knowKinds: knowKinds, clusterName: clusterName, f: f}
	blueprint := actor.NewBlueprint()
	discoveryPid, err := system.Spawn(blueprint, func() actor.IActor { return d }, nil)
	if err != nil {
		zlog.Error("consul new discovery agent err", zap.Error(err))
		return err
	}
	c.discoveryPid = discoveryPid
	zlog.Info("consul new discovery agent", zap.String("clusterName", clusterName), zap.Strings("knowKinds", knowKinds))
	return nil
}

func (c *consulProvider) Register(system actor.ISystem, service common.IMember, ttl, deregisterTtl time.Duration) error {
	// 注册服务
	check := &api.AgentServiceCheck{
		TTL:                            (ttl).String(),
		DeregisterCriticalServiceAfter: (deregisterTtl).String(),
	}
	registration := &api.AgentServiceRegistration{
		ID:      service.GetID(),
		Name:    service.GetName(),
		Address: service.GetAddress(),
		Port:    service.GetPort(),
		Tags:    service.GetTags(),
		Meta:    service.GetMeta(),
		Check:   check,
	}
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		zlog.Error("consul service  register err", zap.String("serviceId", service.GetID()), zap.Error(err))
		return err
	}

	// 启动健康检测actor
	if err := c.spawnHealthCheckActor(system, service, ttl); err != nil {
		zlog.Error("consul new health  agent err", zap.Error(err))
		return err
	}

	zlog.Info("consul service  register ", zap.String("serviceId", service.GetID()),
		zap.String("serviceName", service.GetName()), zap.String("address", service.GetAddress()),
		zap.Int("port", service.GetPort()), zap.Strings("tags", service.GetTags()))
	return nil
}

func (c *consulProvider) UpdateStatus(status string) error {
	return c.healthPid.Send("UpdateStatus", status)
}

func (c *consulProvider) spawnHealthCheckActor(system actor.ISystem, service common.IMember, ttl time.Duration) error {
	if c.healthPid != nil {
		return nil
	}
	h := &healthCheckActor{client: c.client, nodeId: service.GetID(), ttl: ttl, status: "passing"}
	blueprint := actor.NewBlueprint()
	healthPid, err := system.Spawn(blueprint, func() actor.IActor { return h }, nil)
	if err != nil {
		zlog.Error("consul new health  agent err", zap.Error(err))
		return err
	}
	c.healthPid = healthPid
	zlog.Info("consul new health  agent", zap.String("serviceId", service.GetID()))
	return nil
}

func (c *consulProvider) Deregister(serviceID string) error {
	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		zlog.Error("consul service deregister", zap.String("serviceId", serviceID), zap.Error(err))
		return err
	}
	if err := c.healthPid.Stop(); err != nil {
		zlog.Error("consul service deregister", zap.String("serviceId", serviceID), zap.Error(err))
		return err
	}
	c.healthPid = nil
	zlog.Info("consul service deregister", zap.String("serviceId", serviceID))
	return nil
}

func (c *consulProvider) Stop() {
	if c.healthPid != nil {
		_ = c.healthPid.Stop()
	}
	if c.discoveryPid != nil {
		_ = c.discoveryPid.Stop()
	}
	c.discoveryPid = nil
	c.healthPid = nil
}

func (c *consulProvider) PutKV(key string, value []byte) error {
	kv := c.client.KV()
	p := &api.KVPair{Key: key, Value: value}
	_, err := kv.Put(p, nil)
	return err
}

func (c *consulProvider) GetKV(key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}

func (c *consulProvider) DeleteKV(key string) error {
	if _, err := c.client.KV().Delete(key, nil); err != nil {
		return err
	}
	return nil
}

// WatchKV
// @Description: 并发不安全
// @receiver c
// @param key
// @return []string
// @return error
func (c *consulProvider) WatchKV(key string) ([]string, error) {
	q := &api.QueryOptions{WaitIndex: c.kvWaitIndex}
	v, meta, err := c.client.KV().Keys(key, "", q)
	if err != nil {
		return nil, err
	}
	c.kvWaitIndex = meta.LastIndex
	return v, err
}
