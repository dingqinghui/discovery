/**
 * @Author: dingQingHui
 * @Description:
 * @File: cluster
 * @Version: 1.0.0
 * @Date: 2024/5/6 10:50
 */

package cluster

import (
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"github.com/dingqinghui/zlog"
	"golang.org/x/exp/slices"
)

type memberUpdateMsg struct {
	waitIndex  uint64
	memberDict map[string]common.IMember
}

type cacheActor struct {
	actor.BuiltinActor
	opts       *Options
	memberList *common.MemberList
	topology   *common.Topology
}

func (c *cacheActor) Init(ctx actor.IContext, msg interface{}) {
	c.memberList = common.NewMemberList()
	opts := c.opts
	callback := func(waitIndex uint64, memberDict map[string]common.IMember) {
		_ = ctx.Process().Send("EventMemberUpdateHandler", &memberUpdateMsg{
			waitIndex:  waitIndex,
			memberDict: memberDict,
		})
	}
	if err := opts.provider.Discovery(opts.system, opts.service.GetName(), opts.knowKinds, callback); err != nil {
		return
	}
	if err := opts.provider.Register(opts.system, opts.service, opts.ttl, opts.deregisterTtl); err != nil {
		return
	}
}

func (c *cacheActor) EventMemberUpdateHandler(ctx actor.IContext, msg interface{}) {
	c.memberList = common.NewMemberList()
	_msg := msg.(*memberUpdateMsg)
	c.memberList.UpdateClusterTopology(_msg.memberDict, _msg.waitIndex)
}

func (c *cacheActor) GetByKind(ctx actor.IContext, msg interface{}) {
	var ret []common.IMember
	kind := msg.(string)
	for _, member := range c.memberList.Members {
		if slices.Contains(member.GetTags(), kind) {
			ret = append(ret, member)
		}
	}
	_ = ctx.Respond(ret)
}

func (c *cacheActor) GetById(ctx actor.IContext, msg interface{}) {
	var ret common.IMember
	id := msg.(string)
	for _, member := range c.memberList.Members {
		if member.GetID() == id {
			ret = member
			break
		}
	}
	_ = ctx.Respond(ret)
}

func (c *cacheActor) GetAll(ctx actor.IContext, msg interface{}) {
	var ret []common.IMember
	for _, member := range c.memberList.Members {
		ret = append(ret, member)
	}
	_ = ctx.Respond(ret)
}

func (c *cacheActor) Stop(ctx actor.IContext, msg interface{}) {
	c.opts.provider.Stop()
	zlog.Info("cache actor stop")
}
