/**
 * @Author: dingQingHui
 * @Description:
 * @File: discovery_actor
 * @Version: 1.0.0
 * @Date: 2024/11/15 16:01
 */

package consul

import (
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"github.com/dingqinghui/zlog"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"time"
)

type discoveryActor struct {
	actor.BuiltinActor
	client      *api.Client
	waitIndex   uint64
	clusterName string
	knowKinds   []string
	f           common.EventMemberUpdateHandler
}

func (d *discoveryActor) Init(ctx actor.IContext, msg interface{}) {
	d.MonitorMemberStatusChanges(ctx, msg)
}

func (d *discoveryActor) MonitorMemberStatusChanges(ctx actor.IContext, msg interface{}) {
	ctx.AddTimer(time.Millisecond*1, "MonitorMemberStatusChanges")

	opt := &api.QueryOptions{
		WaitIndex: d.waitIndex,
		WaitTime:  time.Second * 3,
	}
	services, meta, err := d.client.Health().Service(d.clusterName, "", true, opt)
	if err != nil {
		zlog.Error("consul discovery agent err", zap.String("clusterName", d.clusterName), zap.Error(err))
		return
	}
	memberDict := make(map[string]common.IMember)
	for _, service := range services {
		for _, tag := range service.Service.Tags {
			if !slices.Contains(d.knowKinds, tag) {
				continue
			}
			member := &common.BuiltinMember{
				Id:      service.Service.ID,
				Name:    service.Service.Service,
				Address: service.Service.Address,
				Port:    service.Service.Port,
				Tags:    service.Service.Tags,
				Meta:    service.Service.Meta,
			}
			memberDict[member.Id] = member
		}
	}
	d.waitIndex = meta.LastIndex
	d.f(d.waitIndex, memberDict)
	return
}

func (d *discoveryActor) Stop(ctx actor.IContext, msg interface{}) {
	zlog.Info("discoveryActor  stop")
}
