/**
 * @Author: dingQingHui
 * @Description:
 * @File: member
 * @Version: 1.0.0
 * @Date: 2024/11/15 17:22
 */

package common

type BuiltinMember struct {
	Name, Id, Address string
	Port              int
	Tags              []string
	Meta              map[string]string
}

func (b *BuiltinMember) GetName() string {
	return b.Name
}

func (b *BuiltinMember) GetID() string {
	return b.Id
}

func (b *BuiltinMember) GetAddress() string {
	return b.Address
}

func (b *BuiltinMember) GetPort() int {
	return b.Port
}

func (b *BuiltinMember) GetTags() []string {
	return b.Tags
}

func (b *BuiltinMember) GetMeta() map[string]string {
	return b.Meta
}
