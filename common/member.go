/**
 * @Author: dingQingHui
 * @Description:
 * @File: member
 * @Version: 1.0.0
 * @Date: 2024/11/15 17:22
 */

package common

type Member struct {
	Name, Id, Address string
	Port              int
	Tags              []string
	Meta              map[string]string
}

func (b *Member) GetName() string {
	return b.Name
}

func (b *Member) GetID() string {
	return b.Id
}

func (b *Member) GetAddress() string {
	return b.Address
}

func (b *Member) GetPort() int {
	return b.Port
}

func (b *Member) GetTags() []string {
	return b.Tags
}

func (b *Member) GetMeta() map[string]string {
	return b.Meta
}
