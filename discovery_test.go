/**
 * @Author: dingQingHui
 * @Description:
 * @File: discovery_test
 * @Version: 1.0.0
 * @Date: 2024/11/18 11:04
 */

package cluster

import (
	"testing"
	"time"
)

func TestDiscovery(t *testing.T) {
	d := New(WithKnowKinds([]string{"Tags"}))
	d.Init()
	members, err := d.GetByKind("Tags")
	_, _ = members, err
	d.Stop()

	time.Sleep(time.Second * 10)
	println(members, err)
}
