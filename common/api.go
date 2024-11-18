/**
 * @Author: dingQingHui
 * @Description:
 * @File: config
 * @Version: 1.0.0
 * @Date: 2024/5/6 10:48
 */

package common

import (
	"github.com/dingqinghui/actor"
	"time"
)

type IMember interface {
	GetName() string
	GetID() string
	GetAddress() string
	GetPort() int
	GetTags() []string
	GetMeta() map[string]string
}

type IProvider interface {
	Name() string
	Discovery(system actor.ISystem, clusterName string, knowKinds []string, f EventMemberUpdateHandler) error
	Register(system actor.ISystem, service IMember, ttl, deregisterTtl time.Duration) error
	UpdateStatus(status string) error
	Deregister(serviceID string) error
	Stop()
	PutKV(key string, value []byte) error
	GetKV(key string) ([]byte, error)
	DeleteKV(key string) error
	WatchKV(key string) ([]string, error)
}

type EventMemberUpdateHandler func(waitIndex uint64, memberDict map[string]IMember)
