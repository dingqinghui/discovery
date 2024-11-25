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

type hubActor struct {
	actor.BuiltinActor
	opts       *Options
	memberList *common.MemberList
	topology   *common.Topology
}

func (c *hubActor) Init(ctx actor.IContext, msg interface{}) error {
	c.memberList = common.NewMemberList()
	opts := c.opts
	callback := func(waitIndex uint64, memberDict map[string]*common.Member) {
		_ = ctx.Process().Send("EventMemberUpdateHandler", waitIndex, memberDict)
	}
	if opts.service == nil {
		return nil
	}
	if err := opts.provider.Discovery(opts.system, opts.service.GetName(), opts.knowKinds, callback); err != nil {
		return err
	}
	if err := opts.provider.Register(opts.system, opts.service, opts.ttl, opts.deregisterTtl); err != nil {
		return err
	}
	return nil
}

func (c *hubActor) EventMemberUpdateHandler(ctx actor.IContext, waitIndex uint64, memberDict map[string]*common.Member) error {
	c.memberList = common.NewMemberList()
	c.memberList.UpdateClusterTopology(memberDict, waitIndex)
	return nil
}

func (c *hubActor) GetByKind(ctx actor.IContext, kind string, reply *MemberList) error {
	for _, member := range c.memberList.Members {
		if slices.Contains(member.GetTags(), kind) {
			reply.Members = append(reply.Members, member)
		}
	}
	return nil
}

func (c *hubActor) GetById(ctx actor.IContext, nodeId string, reply *common.Member) error {
	for _, member := range c.memberList.Members {
		if member.GetID() == nodeId {
			reply = member
			break
		}
	}
	return nil
}

func (c *hubActor) GetAll(ctx actor.IContext, reply *MemberList) error {
	for _, member := range c.memberList.Members {
		reply.Members = append(reply.Members, member)
	}
	return nil
}

func (c *hubActor) Stop(ctx actor.IContext) error {
	c.opts.provider.Stop()
	zlog.Info("cache actor stop")
	return nil
}
