/**
 * @Author: dingQingHui
 * @Description:
 * @File: consul_provider_test
 * @Version: 1.0.0
 * @Date: 2024/4/25 14:28
 */

package consul

import (
	"fmt"
	"github.com/dingqinghui/actor"
	"github.com/dingqinghui/discovery/common"
	"testing"
	"time"
)

func TestConsulProvider_Register(t *testing.T) {
	provider, _ := NewConsulProvider("127.0.0.1:8500")

	service := &common.Member{
		Name:    "clusterName",
		Id:      "NodeId",
		Address: "127.0.0.1",
		Port:    8081,
		Tags:    []string{"NodeType"},
		Meta:    map[string]string{"meta1": "data1"},
	}
	if err := provider.Register(actor.NewSystem(), service, time.Second*1, time.Second*30); err != nil {
		fmt.Printf("consul register err:%v\n", err)
		return
	}

	time.Sleep(time.Second * 100)
}

func EventMemberUpdateHandler(waitIndex uint64, memberDict map[string]*common.Member) {
	memberList := common.NewMemberList()
	topology := memberList.UpdateClusterTopology(memberDict, waitIndex)
	_ = topology
	println()
}

func TestConsulProvider_Discovery(t *testing.T) {
	provider, _ := NewConsulProvider("127.0.0.1:8500")

	if err := provider.Discovery(actor.NewSystem(), "clusterName", []string{"NodeType"}, EventMemberUpdateHandler); err != nil {
		fmt.Printf("consul register err:%v\n", err)
		return
	}

	time.Sleep(time.Second * 10)
	provider.Stop()
	time.Sleep(time.Second * 10)
}
