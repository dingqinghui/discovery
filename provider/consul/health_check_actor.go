/**
 * @Author: dingQingHui
 * @Description:
 * @File: health_check_actor
 * @Version: 1.0.0
 * @Date: 2024/11/15 16:02
 */

package consul

import (
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/zlog"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"time"
)

type healthCheckActor struct {
	actor.BuiltinActor
	client *api.Client
	nodeId string
	status string
	ttl    time.Duration
}

func (h *healthCheckActor) Init(ctx actor.IContext, msg interface{}) error {
	ctx.AddTimer(h.ttl/3, "OnTimer")
	return nil
}

func (h *healthCheckActor) UpdateStatus(ctx actor.IContext, status string) error {
	h.status = status
	h.updateStatus()
	return nil
}

func (h *healthCheckActor) OnTimer(ctx actor.IContext) error {
	h.updateStatus()
	ctx.AddTimer(h.ttl/3, "OnTimer")
	return nil
}

func (h *healthCheckActor) updateStatus() {
	if err := h.client.Agent().UpdateTTL("service:"+h.nodeId, "", h.status); err != nil {
		zlog.Error("consul health agent err", zap.String("serviceId", h.nodeId),
			zap.String("status", h.status), zap.Error(err))
		return
	}
}

func (h *healthCheckActor) Stop(ctx actor.IContext) error {
	zlog.Info("healthCheckActor  stop")
	return nil
}
